package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

const unreadableResponse = "UNREADABLE"
const canvasAnalyzeAttempts = 3

func (g *gateway) AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error) {
	if !g.IsConfigured(cfg) {
		return "", errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	prompt := buildCanvasPrompt(correctAnswer)
	log.Printf("[assistant] analyze_canvas prompt=%q", truncateForLog(prompt, 700))

	var lastErr error
	for attempt := 1; attempt <= canvasAnalyzeAttempts; attempt++ {
		log.Printf("[assistant] analyze_canvas attempt=%d/%d host=%s", attempt, canvasAnalyzeAttempts, baseURL)
		conversationID, err := g.createConversation(ctx, baseURL, apiKey)
		if err != nil {
			log.Printf("[assistant] analyze_canvas create_conversation_error attempt=%d err=%v", attempt, err)
			lastErr = err
			continue
		}
		log.Printf("[assistant] analyze_canvas conversation_created attempt=%d conversation_id=%s", attempt, conversationID)

		response, err := g.sendCanvasMessage(ctx, baseURL, apiKey, conversationID, prompt, canvasData)
		if err != nil {
			log.Printf("[assistant] analyze_canvas send_canvas_error attempt=%d conversation_id=%s err=%v", attempt, conversationID, err)
			lastErr = err
			continue
		}
		if !isExpectedCanvasResponse(response) {
			log.Printf("[assistant] analyze_canvas unexpected_response attempt=%d conversation_id=%s response=%q", attempt, conversationID, response)
			lastErr = errors.New("assistant canvas response format not expected")
			continue
		}
		normalized := normalizeCanvasResponse(response)
		log.Printf("[assistant] analyze_canvas success attempt=%d conversation_id=%s response=%q normalized=%q", attempt, conversationID, response, normalized)
		return normalized, nil
	}

	if lastErr == nil {
		lastErr = errors.New("assistant canvas analysis failed")
	}
	return "", fmt.Errorf("assistant canvas analysis failed after %d attempts: %w", canvasAnalyzeAttempts, lastErr)
}

func (g *gateway) sendCanvasMessage(ctx context.Context, baseURL, apiKey, conversationID, prompt, canvasData string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("content", prompt); err != nil {
		return "", err
	}

	decoded, contentType, err := utils.DecodeDataURI(canvasData)
	if err != nil {
		return "", err
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="image_content"; filename="canvas.png"`)
	partHeader.Set("Content-Type", contentType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(decoded); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/conversation/"+conversationID+"/message",
		&body,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("[assistant] send_canvas_message request_error host=%s conversation_id=%s err=%v", baseURL, conversationID, err)
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[assistant] send_canvas_message status=%d host=%s conversation_id=%s", resp.StatusCode, baseURL, conversationID)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("assistant add message returned status %d", resp.StatusCode)
	}

	var parsed messageResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	for i := len(parsed.Data) - 1; i >= 0; i-- {
		msg := parsed.Data[i]
		if msg.Sender == "assistant" && strings.TrimSpace(msg.Content) != "" {
			assistantReply := strings.TrimSpace(msg.Content)
			if strings.EqualFold(normalizeCanvasResponse(assistantReply), unreadableResponse) {
				if fallback := extractTextFound(parsed.Data); fallback != "" {
					return fallback, nil
				}
			}
			return assistantReply, nil
		}
	}

	return "", errors.New("assistant response missing")
}

func buildCanvasPrompt(correctAnswer string) string {
	cleanAnswer := strings.TrimSpace(correctAnswer)
	if cleanAnswer == "" {
		cleanAnswer = "(sin respuesta correcta provista)"
	}

	return "Analiza esta imagen de una respuesta manuscrita de un estudiante. " +
		"Extrae únicamente la respuesta final escrita. " +
		"No expliques el procedimiento. " +
		"Si no puedes leerla con suficiente confianza, responde exactamente: " + unreadableResponse + ". " +
		"Respuesta correcta esperada: " + cleanAnswer + "."
}

func isExpectedCanvasResponse(raw string) bool {
	value := normalizeCanvasResponse(raw)
	if value == "" {
		return false
	}
	if strings.EqualFold(value, unreadableResponse) {
		return true
	}
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") || strings.Contains(value, "```") {
		return false
	}
	if len(value) > 120 {
		return false
	}
	return true
}

func normalizeCanvasResponse(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}

	lower := strings.ToLower(value)
	if strings.Contains(lower, unreadableResponse) {
		return unreadableResponse
	}

	if strings.Contains(value, "</think>") {
		parts := strings.Split(value, "</think>")
		value = strings.TrimSpace(parts[len(parts)-1])
	}

	value = strings.ReplaceAll(value, "```", "")
	value = strings.TrimSpace(value)
	if strings.Contains(strings.ToLower(value), unreadableResponse) {
		return unreadableResponse
	}
	return value
}

func extractTextFound(messages []struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}) string {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := strings.TrimSpace(messages[i].Content)
		if msg == "" {
			continue
		}
		upper := strings.ToUpper(msg)
		idx := strings.LastIndex(upper, "TEXT_FOUND:")
		if idx == -1 {
			continue
		}
		segment := msg[idx+len("TEXT_FOUND:"):]
		end := len(segment)
		for _, marker := range []string{"\n", "KEY_FINDINGS:", "RESPONSE_CONTEXT:"} {
			markerIdx := strings.Index(strings.ToUpper(segment), marker)
			if markerIdx != -1 && markerIdx < end {
				end = markerIdx
			}
		}
		value := strings.TrimSpace(segment[:end])
		value = strings.Trim(value, "-: \"'")
		if value == "" || strings.EqualFold(value, "none") || strings.EqualFold(value, unreadableResponse) {
			continue
		}
		if len(value) > 32 {
			continue
		}
		return value
	}
	return ""
}

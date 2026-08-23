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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

const unreadableResponse = "UNREADABLE"
const canvasAnalyzeAttempts = 3

func (g *gateway) AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error) {
	return g.analyzeCanvas(ctx, cfg, canvasData, buildCanvasPrompt(correctAnswer), isExpectedCanvasResponse)
}

// AnalyzeNotebookCanvas reads a whole notebook page rather than one exercise's
// final answer. It cannot share AnalyzeCanvas's prompt: that one asks for the
// final answer only and rejects anything past 120 characters, so a faithful
// transcription of a page of worked problems was thrown away as malformed and
// retried until the attempts ran out.
func (g *gateway) AnalyzeNotebookCanvas(ctx context.Context, cfg Config, canvasData, pageContext string) (string, error) {
	return g.analyzeCanvas(ctx, cfg, canvasData, buildNotebookCanvasPrompt(pageContext), isExpectedNotebookResponse)
}

func (g *gateway) AnalyzeNotebookStatement(ctx context.Context, cfg Config, imageData, pageContext string) (string, error) {
	return g.analyzeCanvas(ctx, cfg, imageData, buildNotebookStatementPrompt(pageContext), isExpectedNotebookResponse)
}

func (g *gateway) analyzeCanvas(
	ctx context.Context,
	cfg Config,
	canvasData, prompt string,
	accept func(string) bool,
) (string, error) {
	if !g.IsConfigured(cfg) {
		return "", errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

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
		if !accept(response) {
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
	dumpCanvasForDebug(decoded, contentType)

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

	responseBody, err := readAssistantResponseBody(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed messageResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
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

func dumpCanvasForDebug(content []byte, contentType string) {
	dir := strings.TrimSpace(os.Getenv("PRACTIQ_DEBUG_CANVAS_DIR"))
	if dir == "" {
		return
	}
	ext := ".png"
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		ext = ".jpg"
	}
	name := filepath.Join(dir, fmt.Sprintf("canvas-%d%s", time.Now().UnixNano(), ext))
	if err := os.WriteFile(name, content, 0o600); err != nil {
		log.Printf("[assistant] canvas debug dump failed err=%v", err)
		return
	}
	log.Printf("[assistant] canvas debug dump path=%s bytes=%d content_type=%s", name, len(content), contentType)
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

func buildNotebookCanvasPrompt(pageContext string) string {
	context := strings.TrimSpace(pageContext)
	if context == "" {
		context = "(sin contexto de la pagina)"
	}

	return "Analiza esta imagen de una pagina de cuaderno resuelta a mano por un estudiante. " +
		"Transcribi todo lo que el estudiante escribio, respetando el orden y separando cada ejercicio con un salto de linea. " +
		"No corrijas, no resuelvas y no opines sobre si esta bien o mal: solo transcribi. " +
		"Si la pagina esta vacia o no podes leer nada con suficiente confianza, responde exactamente: " + unreadableResponse + ". " +
		"Contexto de la pagina: " + context + "."
}

func buildNotebookStatementPrompt(pageContext string) string {
	context := strings.TrimSpace(pageContext)
	if context == "" {
		context = "(sin contexto de la pagina)"
	}

	return "Analiza esta imagen de una pagina de cuaderno preparada por un docente. " +
		"Transcribi la consigna tal como esta escrita: los enunciados, los ejercicios y, si estan, las respuestas esperadas. " +
		"Separa cada ejercicio con un salto de linea. No la resuelvas ni agregues nada que no este en la imagen. " +
		"Si la pagina esta vacia o no podes leerla con suficiente confianza, responde exactamente: " + unreadableResponse + ". " +
		"Contexto de la pagina: " + context + "."
}

// A page of worked problems is legitimately long, so this only guards against
// the assistant answering with something that is not a transcription at all.
func isExpectedNotebookResponse(raw string) bool {
	value := normalizeCanvasResponse(raw)
	if value == "" {
		return false
	}
	if strings.EqualFold(value, unreadableResponse) {
		return true
	}
	if isCanvasVerdict(value) {
		return false
	}
	// No fence check here, unlike the single-answer version below:
	// normalizeCanvasResponse strips ``` before this runs, so testing for it
	// can never fire. A fenced answer arrives already unwrapped and usable.
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return false
	}
	return len(value) <= maxNotebookTranscriptionChars
}

// Long enough for a full page, short enough that a runaway answer still fails.
const maxNotebookTranscriptionChars = 4000

func isExpectedCanvasResponse(raw string) bool {
	value := normalizeCanvasResponse(raw)
	if value == "" {
		return false
	}
	if strings.EqualFold(value, unreadableResponse) {
		return true
	}
	// Gillie's conversational tutor can return a review verdict instead of
	// transcribing the handwriting. A verdict is never a student answer.
	if isCanvasVerdict(value) {
		return false
	}
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") || strings.Contains(value, "```") {
		return false
	}
	if len(value) > 120 {
		return false
	}
	return true
}

func isCanvasVerdict(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.Trim(normalized, ".,!¡!¿?;: \t\n\r")
	switch normalized {
	case "correcta", "correcto", "incorrecta", "incorrecto", "es correcta", "es correcto", "es incorrecta", "es incorrecto":
		return true
	default:
		return false
	}
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

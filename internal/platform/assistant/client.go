package assistant

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

const unreadableResponse = "UNREADABLE"
const canvasAnalyzeAttempts = 3

type Config struct {
	BaseURL string
	APIKey  string
}

type Service interface {
	IsConfigured(cfg Config) bool
	AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error)
	EvaluatePracticeAnswer(ctx context.Context, cfg Config, question, correctAnswer, studentAnswer string) (EvaluationResult, error)
}

type EvaluationResult struct {
	IsCorrect bool
	Feedback  string
}

type service struct {
	client *http.Client
}

type createConversationRequest struct {
	Title     string `json:"title"`
	IsSandbox bool   `json:"is_sandbox"`
}

type createConversationResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type messageResponse struct {
	Data []struct {
		Content string `json:"content"`
		Sender  string `json:"sender"`
	} `json:"data"`
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *service) IsConfigured(cfg Config) bool {
	return strings.TrimSpace(cfg.BaseURL) != "" && strings.TrimSpace(cfg.APIKey) != ""
}

func (s *service) AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error) {
	if !s.IsConfigured(cfg) {
		return "", errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	prompt := buildPrompt(correctAnswer)
	log.Printf("[assistant] analyze_canvas prompt=%q", truncateForLog(prompt, 700))

	var lastErr error
	for attempt := 1; attempt <= canvasAnalyzeAttempts; attempt++ {
		log.Printf("[assistant] analyze_canvas attempt=%d/%d host=%s", attempt, canvasAnalyzeAttempts, baseURL)
		conversationID, err := s.createConversation(ctx, baseURL, apiKey)
		if err != nil {
			log.Printf("[assistant] analyze_canvas create_conversation_error attempt=%d err=%v", attempt, err)
			lastErr = err
			continue
		}
		log.Printf("[assistant] analyze_canvas conversation_created attempt=%d conversation_id=%s", attempt, conversationID)

		response, err := s.sendCanvasMessage(ctx, baseURL, apiKey, conversationID, prompt, canvasData)
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

func (s *service) EvaluatePracticeAnswer(ctx context.Context, cfg Config, question, correctAnswer, studentAnswer string) (EvaluationResult, error) {
	if !s.IsConfigured(cfg) {
		return EvaluationResult{}, errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	conversationID, err := s.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		log.Printf("[assistant] evaluate_answer create_conversation_error host=%s err=%v", baseURL, err)
		return EvaluationResult{}, err
	}
	log.Printf("[assistant] evaluate_answer conversation_created host=%s conversation_id=%s", baseURL, conversationID)

	prompt := buildEvaluationPrompt(question, correctAnswer, studentAnswer)
	log.Printf("[assistant] evaluate_answer prompt=%q", truncateForLog(prompt, 700))
	response, err := s.sendTextMessage(ctx, baseURL, apiKey, conversationID, prompt)
	if err != nil {
		log.Printf("[assistant] evaluate_answer send_text_error conversation_id=%s err=%v", conversationID, err)
		return EvaluationResult{}, err
	}
	log.Printf("[assistant] evaluate_answer raw_response conversation_id=%s response=%q", conversationID, response)

	evaluation, parseErr := parseEvaluationResponse(response)
	if parseErr != nil {
		log.Printf("[assistant] evaluate_answer parse_error conversation_id=%s err=%v", conversationID, parseErr)
		return EvaluationResult{}, parseErr
	}
	log.Printf("[assistant] evaluate_answer success conversation_id=%s is_correct=%t feedback=%q", conversationID, evaluation.IsCorrect, evaluation.Feedback)
	return evaluation, nil
}

func (s *service) createConversation(ctx context.Context, baseURL, apiKey string) (string, error) {
	body, err := json.Marshal(createConversationRequest{
		Title:     "Practiq Canvas Evaluation",
		IsSandbox: false,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/conversation/", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[assistant] create_conversation request_error host=%s err=%v", baseURL, err)
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[assistant] create_conversation status=%d host=%s", resp.StatusCode, baseURL)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("assistant create conversation returned status %d", resp.StatusCode)
	}

	var parsed createConversationResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if parsed.Data.ID == "" {
		return "", errors.New("assistant conversation id missing")
	}

	return parsed.Data.ID, nil
}

func (s *service) sendCanvasMessage(ctx context.Context, baseURL, apiKey, conversationID, prompt, canvasData string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("content", prompt); err != nil {
		return "", err
	}

	decoded, contentType, err := decodeDataURI(canvasData)
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

	resp, err := s.client.Do(req)
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

func (s *service) sendTextMessage(ctx context.Context, baseURL, apiKey, conversationID, prompt string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("content", prompt); err != nil {
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

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[assistant] send_text_message request_error host=%s conversation_id=%s err=%v", baseURL, conversationID, err)
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[assistant] send_text_message status=%d host=%s conversation_id=%s", resp.StatusCode, baseURL, conversationID)

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
			return strings.TrimSpace(msg.Content), nil
		}
	}

	return "", errors.New("assistant response missing")
}

func decodeDataURI(dataURI string) ([]byte, string, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid canvas data uri")
	}

	meta := parts[0]
	payload := parts[1]
	contentType := "image/png"

	if strings.HasPrefix(meta, "data:") {
		semi := strings.Index(meta, ";")
		if semi > len("data:") {
			contentType = meta[len("data:"):semi]
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", err
	}

	return decoded, contentType, nil
}

func buildPrompt(correctAnswer string) string {
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

func buildEvaluationPrompt(question, correctAnswer, studentAnswer string) string {
	q := strings.TrimSpace(question)
	if q == "" {
		q = "(sin enunciado provisto)"
	}
	ca := strings.TrimSpace(correctAnswer)
	if ca == "" {
		ca = "(sin respuesta correcta provista)"
	}
	sa := strings.TrimSpace(studentAnswer)
	if sa == "" {
		sa = "(sin respuesta del estudiante)"
	}

	return "Evalua si la respuesta del estudiante es correcta para el ejercicio. " +
		"Responde SOLO con JSON valido en una sola linea con este formato exacto: " +
		`{"is_correct":true,"feedback":"comentario breve"}` + ". " +
		"No incluyas texto adicional. " +
		"El campo feedback debe ser breve y explicar por que. " +
		"Enunciado: " + q + ". " +
		"Respuesta correcta esperada: " + ca + ". " +
		"Respuesta del estudiante: " + sa + "."
}

func parseEvaluationResponse(raw string) (EvaluationResult, error) {
	trimmed := strings.TrimSpace(raw)
	body := ""

	if idx := strings.LastIndex(strings.ToLower(trimmed), `{"is_correct"`); idx != -1 {
		candidate := trimmed[idx:]
		if end := strings.Index(candidate, "}"); end != -1 {
			body = candidate[:end+1]
		}
	}

	if body == "" {
		start := strings.Index(trimmed, "{")
		end := strings.LastIndex(trimmed, "}")
		if start == -1 || end == -1 || end < start {
			return EvaluationResult{}, errors.New("assistant evaluation response is not JSON")
		}
		body = trimmed[start : end+1]
	}

	var parsed struct {
		IsCorrect bool   `json:"is_correct"`
		Feedback  string `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return EvaluationResult{}, err
	}

	return EvaluationResult{
		IsCorrect: parsed.IsCorrect,
		Feedback:  normalizeFeedback(parsed.Feedback),
	}, nil
}

func normalizeFeedback(value string) string {
	feedback := strings.TrimSpace(value)
	feedback = strings.ReplaceAll(feedback, "\n", " ")
	feedback = strings.ReplaceAll(feedback, "\t", " ")
	feedback = strings.Join(strings.Fields(feedback), " ")
	if len(feedback) > 500 {
		feedback = strings.TrimSpace(feedback[:500])
	}
	return feedback
}

func truncateForLog(value string, max int) string {
	if max <= 0 || len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
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

package assistant

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

const unreadableResponse = "UNREADABLE"

type Config struct {
	BaseURL string
	APIKey  string
}

type Service interface {
	IsConfigured(cfg Config) bool
	AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error)
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

	conversationID, err := s.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		return "", err
	}

	prompt := buildPrompt(correctAnswer)
	return s.sendCanvasMessage(ctx, baseURL, apiKey, conversationID, prompt, canvasData)
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
		return "", err
	}
	defer resp.Body.Close()

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
	partHeader.Set("Content-Disposition", `form-data; name="image_file"; filename="canvas.png"`)
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
		return "", err
	}
	defer resp.Body.Close()

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

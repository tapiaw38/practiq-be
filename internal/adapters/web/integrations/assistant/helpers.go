package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/urlvalidator"
)

func (g *gateway) createConversation(ctx context.Context, baseURL, apiKey string) (string, error) {
	if err := validateAssistantURL(baseURL); err != nil {
		return "", fmt.Errorf("assistant URL validation failed: %w", err)
	}

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

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("[assistant] create_conversation request_error host=%s err=%v", baseURL, err)
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[assistant] create_conversation status=%d host=%s", resp.StatusCode, baseURL)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("assistant create conversation returned status %d", resp.StatusCode)
	}

	responseBody, err := readAssistantResponseBody(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed createConversationResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return "", err
	}
	if parsed.Data.ID == "" {
		return "", errors.New("assistant conversation id missing")
	}

	return parsed.Data.ID, nil
}

func (g *gateway) sendTextMessage(ctx context.Context, baseURL, apiKey, conversationID, prompt string) (string, error) {
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

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("[assistant] send_text_message request_error host=%s conversation_id=%s err=%v", baseURL, conversationID, err)
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[assistant] send_text_message status=%d host=%s conversation_id=%s", resp.StatusCode, baseURL, conversationID)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("assistant add message returned status %d", resp.StatusCode)
	}

	return lastAssistantMessage(resp.Body)
}

// lastAssistantMessage picks the newest assistant reply out of the conversation
// payload the assistant returns after a message is posted.
func lastAssistantMessage(body io.Reader) (string, error) {
	responseBody, err := readAssistantResponseBody(body)
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
			return strings.TrimSpace(msg.Content), nil
		}
	}

	return "", errors.New("assistant response missing")
}

func validateAssistantURL(rawURL string) error {
	return urlvalidator.ValidateURLWithOptions(rawURL, urlvalidator.Options{
		AllowedPrivateHostnames: allowedPrivateHostnames(),
	})
}

func readAssistantResponseBody(body io.Reader) ([]byte, error) {
	responseBody, err := io.ReadAll(io.LimitReader(body, maxAssistantResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(responseBody) > maxAssistantResponseBytes {
		return nil, fmt.Errorf("assistant response exceeds %d bytes", maxAssistantResponseBytes)
	}
	return responseBody, nil
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

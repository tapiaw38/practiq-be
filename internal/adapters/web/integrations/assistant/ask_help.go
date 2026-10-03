package assistant

import (
	"context"
	"errors"
	"log"
	"strings"
)

func (g *gateway) AskHelp(ctx context.Context, cfg Config, prompt string) (string, error) {
	if !g.IsConfigured(cfg) {
		return "", errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	conversationID, err := g.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		log.Printf("[assistant] ask_help create_conversation_error host=%s err=%v", baseURL, err)
		return "", err
	}
	log.Printf("[assistant] ask_help conversation_created host=%s conversation_id=%s", baseURL, conversationID)

	log.Printf("[assistant] ask_help prompt=%q", truncateForLog(prompt, 700))
	response, err := g.sendTextMessage(ctx, baseURL, apiKey, conversationID, prompt)
	if err != nil {
		log.Printf("[assistant] ask_help send_text_error conversation_id=%s err=%v", conversationID, err)
		return "", err
	}
	log.Printf("[assistant] ask_help success conversation_id=%s response=%q", conversationID, truncateForLog(response, 200))
	return response, nil
}

package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func (g *gateway) GenerateCourseCuriosities(ctx context.Context, cfg Config, subject, topic string, count int) ([]string, error) {
	if !g.IsConfigured(cfg) {
		return nil, errors.New("assistant service not configured")
	}

	if count <= 0 {
		count = 5
	}
	if count > 10 {
		count = 10
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	conversationID, err := g.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		log.Printf("[assistant] generate_curiosities create_conversation_error host=%s err=%v", baseURL, err)
		return nil, err
	}
	log.Printf("[assistant] generate_curiosities conversation_created host=%s conversation_id=%s", baseURL, conversationID)

	prompt := buildCuriositiesPrompt(subject, topic, count)
	log.Printf("[assistant] generate_curiosities prompt=%q", truncateForLog(prompt, 500))

	response, err := g.sendTextMessage(ctx, baseURL, apiKey, conversationID, prompt)
	if err != nil {
		log.Printf("[assistant] generate_curiosities send_text_error conversation_id=%s err=%v", conversationID, err)
		return nil, err
	}
	log.Printf("[assistant] generate_curiosities success conversation_id=%s response=%q", conversationID, truncateForLog(response, 300))

	curiosities, err := parseCuriositiesResponse(response)
	if err != nil {
		log.Printf("[assistant] generate_curiosities parse_error response=%q err=%v", truncateForLog(response, 200), err)
		return nil, err
	}

	return curiosities, nil
}

func buildCuriositiesPrompt(subject, topic string, count int) string {
	return fmt.Sprintf(`Genera exactamente %d datos curiosos cortos y divertidos para niños sobre "%s" en la materia "%s".

Requisitos:
- Cada dato debe tener máximo 100 caracteres
- Usa lenguaje simple y amigable para niños de primaria
- Que sean sorprendentes o interesantes
- Evita datos muy obvios o aburridos

Responde SOLO con un array JSON de strings, sin explicación adicional.
Ejemplo de formato: ["Dato curioso 1", "Dato curioso 2"]`, count, topic, subject)
}

func parseCuriositiesResponse(response string) ([]string, error) {
	response = strings.TrimSpace(response)

	// Try to extract JSON array from response
	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")

	if start == -1 || end == -1 || end <= start {
		// If no JSON array found, try to parse as plain text lines
		return parseAsLines(response), nil
	}

	jsonStr := response[start : end+1]

	var curiosities []string
	if err := json.Unmarshal([]byte(jsonStr), &curiosities); err != nil {
		return parseAsLines(response), nil
	}

	// Filter empty strings and trim
	var result []string
	for _, c := range curiosities {
		c = strings.TrimSpace(c)
		if c != "" {
			result = append(result, c)
		}
	}

	return result, nil
}

func parseAsLines(response string) []string {
	lines := strings.Split(response, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Remove common prefixes like "- ", "1. ", "• "
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "• ")
		if len(line) > 2 && line[1] == '.' {
			line = strings.TrimSpace(line[2:])
		}
		if line != "" && len(line) > 10 {
			result = append(result, line)
		}
	}
	return result
}

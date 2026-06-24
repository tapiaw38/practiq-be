package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

func (g *gateway) EvaluatePracticeAnswer(ctx context.Context, cfg Config, question, correctAnswer, studentAnswer, gradeName string) (EvaluationResult, error) {
	if !g.IsConfigured(cfg) {
		return EvaluationResult{}, errors.New("assistant service not configured")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	conversationID, err := g.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		log.Printf("[assistant] evaluate_answer create_conversation_error host=%s err=%v", baseURL, err)
		return EvaluationResult{}, err
	}
	log.Printf("[assistant] evaluate_answer conversation_created host=%s conversation_id=%s", baseURL, conversationID)

	prompt := buildEvaluationPrompt(question, correctAnswer, studentAnswer, gradeName)
	log.Printf("[assistant] evaluate_answer prompt=%q", truncateForLog(prompt, 700))
	response, err := g.sendTextMessage(ctx, baseURL, apiKey, conversationID, prompt)
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

func buildEvaluationPrompt(question, correctAnswer, studentAnswer, gradeName string) string {
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
	gn := strings.TrimSpace(gradeName)

	gradeContext := ""
	if gn != "" {
		gradeContext = "Considera los contenidos curriculares y documentos de " + gn + " para evaluar. "
	}

	return "Evalua si la respuesta del estudiante es correcta para el ejercicio. " +
		gradeContext +
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

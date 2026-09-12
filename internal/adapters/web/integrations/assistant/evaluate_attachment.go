package assistant

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
)

var multipartQuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeMultipartQuotes(value string) string {
	return multipartQuoteEscaper.Replace(value)
}

// EvaluateAttachment grades a file the student uploaded as their answer. Audio
// and images use media channels; PDF/DOCX use Gillie's document channel.
func (g *gateway) EvaluateAttachment(ctx context.Context, cfg Config, input AttachmentEvaluationInput) (EvaluationResult, error) {
	if !g.IsConfigured(cfg) {
		return EvaluationResult{}, errors.New("assistant service not configured")
	}

	field, err := attachmentFieldFor(input.Kind)
	if err != nil {
		return EvaluationResult{}, err
	}
	if len(input.Content) == 0 {
		return EvaluationResult{}, errors.New("attachment is empty")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	conversationID, err := g.createConversation(ctx, baseURL, apiKey)
	if err != nil {
		log.Printf("[assistant] evaluate_attachment create_conversation_error host=%s err=%v", baseURL, err)
		return EvaluationResult{}, err
	}

	prompt := buildAttachmentEvaluationPrompt(input)
	response, err := g.sendAttachmentMessage(ctx, baseURL, apiKey, conversationID, prompt, field, input.Filename, attachmentContentType(input), input.Content)
	if err != nil {
		log.Printf("[assistant] evaluate_attachment send_error conversation_id=%s err=%v", conversationID, err)
		return EvaluationResult{}, err
	}
	log.Printf("[assistant] evaluate_attachment success conversation_id=%s response=%q", conversationID, truncateForLog(response, 200))

	return parseEvaluationResponse(response)
}

// attachmentFieldFor maps a file kind to the assistant's multipart field.
func attachmentFieldFor(kind string) (string, error) {
	switch kind {
	case AttachmentKindAudio:
		return "voice_content", nil
	case AttachmentKindImage:
		return "image_content", nil
	case AttachmentKindPDF, AttachmentKindDocument:
		return "document_content", nil
	default:
		return "", fmt.Errorf("%w: %s", ErrAttachmentNotEvaluable, kind)
	}
}

func buildAttachmentEvaluationPrompt(input AttachmentEvaluationInput) string {
	var sb strings.Builder

	if input.Kind == AttachmentKindAudio {
		sb.WriteString("El alumno respondió con un audio adjunto. Escuchá el audio y evaluá la respuesta.\n")
	} else if input.Kind == AttachmentKindImage {
		sb.WriteString("El alumno respondió con una imagen adjunta. Analizá la imagen y evaluá la respuesta.\n")
	} else {
		sb.WriteString("El alumno respondió con un documento adjunto. Usá las herramientas documentales disponibles para leerlo y evaluá la respuesta.\n")
		fmt.Fprintf(&sb, "Nombre del documento: %s\n", input.Filename)
	}
	fmt.Fprintf(&sb, "Consigna: %s\n", input.Question)
	if input.CorrectAnswer != "" {
		fmt.Fprintf(&sb, "Respuesta esperada: %s\n", input.CorrectAnswer)
	}
	if input.GradeName != "" {
		fmt.Fprintf(&sb, "El alumno es de %s; evaluá con ese nivel en mente.\n", input.GradeName)
	}
	sb.WriteString(`
Respondé SOLO con un JSON: {"is_correct": true|false, "feedback": "..."}
El feedback es para un niño: breve, claro y alentador.
Si el archivo no se entiende o no se puede evaluar, respondé {"is_correct": false, "feedback": "UNREADABLE"}.`)

	return sb.String()
}

func attachmentContentType(input AttachmentEvaluationInput) string {
	if declared := strings.TrimSpace(input.ContentType); declared != "" {
		return declared
	}
	if byExt := mime.TypeByExtension(strings.ToLower(filepath.Ext(input.Filename))); byExt != "" {
		return byExt
	}
	return "application/octet-stream"
}

func (g *gateway) sendAttachmentMessage(ctx context.Context, baseURL, apiKey, conversationID, prompt, field, filename, contentType string, content []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("content", prompt); err != nil {
		return "", err
	}
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(
		`form-data; name="%s"; filename="%s"`,
		escapeMultipartQuotes(field),
		escapeMultipartQuotes(filename),
	))
	partHeader.Set("Content-Type", contentType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(content); err != nil {
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
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("assistant attachment message returned status %d", resp.StatusCode)
	}
	return lastAssistantMessage(resp.Body)
}

package practicesheet

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

const exerciseTypeAttachment = "attachment"

// attachmentOutcome is what the submit flow needs to know about a file answer.
type attachmentOutcome struct {
	IsCorrect bool
	Feedback  string
	// NeedsReview means nobody graded it yet: the answer is excluded from the
	// score instead of counting as wrong.
	NeedsReview bool
}

// evaluateAttachment uses Gillie's media/document channels and defers to a
// teacher for everything it cannot read. It never marks an answer wrong just
// because grading was impossible.
func evaluateAttachment(
	ctx context.Context,
	app *appcontext.Context,
	cfg assistant.Config,
	ex domain.Exercise,
	gradeName string,
	attachmentURL, filename, contentType string,
) attachmentOutcome {
	pending := attachmentOutcome{NeedsReview: true, Feedback: "Tu entrega quedó pendiente de revisión del docente."}

	kind, _, err := storage.ClassifyContentType(contentType)
	if err != nil {
		return pending
	}
	if kind != storage.FileKindAudio && kind != storage.FileKindImage &&
		kind != storage.FileKindPDF && kind != storage.FileKindDocument {
		return pending
	}
	if app.Integrations.AssistantGateway == nil || !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return pending
	}
	if app.ImageStorage == nil {
		return pending
	}

	content, _, err := app.ImageStorage.FetchFile(ctx, attachmentURL)
	if err != nil {
		log.Printf("[practice_attachment] could not fetch file url=%q err=%v", attachmentURL, err)
		return pending
	}

	evaluation, err := app.Integrations.AssistantGateway.EvaluateAttachment(ctx, cfg, assistant.AttachmentEvaluationInput{
		Question:      ex.Question,
		CorrectAnswer: ex.CorrectAnswer,
		GradeName:     gradeName,
		Kind:          string(kind),
		Filename:      filename,
		Content:       content,
	})
	if err != nil {
		if !errors.Is(err, assistant.ErrAttachmentNotEvaluable) {
			log.Printf("[practice_attachment] evaluation failed exercise_id=%s err=%v", ex.ID, err)
		}
		return pending
	}
	// The assistant received the file but could not make sense of it.
	if strings.Contains(strings.ToUpper(evaluation.Feedback), assistant.UnreadableFeedback) {
		return pending
	}

	return attachmentOutcome{IsCorrect: evaluation.IsCorrect, Feedback: evaluation.Feedback}
}

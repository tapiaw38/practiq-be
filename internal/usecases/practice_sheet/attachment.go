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

const attachmentsFolder = "attachments"

type attachmentOutcome struct {
	IsCorrect bool
	Feedback  string

	Ungraded bool

	AISuggestedCorrect *bool
}

func evaluateAttachment(
	ctx context.Context,
	app *appcontext.Context,
	cfg assistant.Config,
	ex domain.Exercise,
	gradeName string,
	attachmentURL, filename string,
	teacherGrades bool,
) attachmentOutcome {
	pending := attachmentOutcome{Ungraded: true, Feedback: ungradedAttachmentFeedback(teacherGrades)}

	if app.Integrations.AssistantGateway == nil || !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return pending
	}
	if app.ImageStorage == nil {
		return pending
	}

	content, storedContentType, err := app.ImageStorage.FetchFile(ctx, attachmentURL)
	if err != nil {
		log.Printf("[practice_attachment] could not fetch file url=%q err=%v", attachmentURL, err)
		return pending
	}

	kind, _, err := storage.ClassifyContentType(storedContentType)
	if err != nil || (kind != storage.FileKindAudio && kind != storage.FileKindImage &&
		kind != storage.FileKindPDF && kind != storage.FileKindDocument) {
		return pending
	}
	if !attachmentKindAccepted(ex, string(kind)) {
		log.Printf("[practice_attachment] kind not accepted exercise_id=%s kind=%s", ex.ID, kind)
		return pending
	}

	evaluation, err := app.Integrations.AssistantGateway.EvaluateAttachment(ctx, cfg, assistant.AttachmentEvaluationInput{
		Question:      ex.Question,
		CorrectAnswer: ex.CorrectAnswer,
		GradeName:     gradeName,
		Kind:          string(kind),
		Filename:      filename,
		ContentType:   storedContentType,
		Content:       content,
	})
	if err != nil {
		if !errors.Is(err, assistant.ErrAttachmentNotEvaluable) {
			log.Printf("[practice_attachment] evaluation failed exercise_id=%s err=%v", ex.ID, err)
		}
		return pending
	}

	if strings.Contains(strings.ToUpper(evaluation.Feedback), assistant.UnreadableFeedback) {
		return pending
	}

	verdict := evaluation.IsCorrect
	return attachmentOutcome{
		IsCorrect: verdict,
		Feedback:  evaluation.Feedback,

		Ungraded:           false,
		AISuggestedCorrect: &verdict,
	}
}

func ungradedAttachmentFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "Tu entrega quedó pendiente de revisión del docente."
	}
	return "No pudimos corregir esta entrega automáticamente, así que no cuenta en tu puntaje."
}

func attachmentKindAccepted(ex domain.Exercise, kind string) bool {
	accepted := ex.AcceptedAttachmentKinds()
	if len(accepted) == 0 {
		return true
	}
	for _, allowed := range accepted {
		if strings.EqualFold(strings.TrimSpace(allowed), kind) {
			return true
		}
	}
	return false
}

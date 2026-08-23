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

// attachmentsFolder is the bucket prefix student deliveries are uploaded under;
// it is what ownership is checked against.
const attachmentsFolder = "attachments"

// attachmentOutcome is what the submit flow needs to know about a file answer.
type attachmentOutcome struct {
	IsCorrect bool
	Feedback  string
	// NeedsReview means nobody could grade it: the answer is excluded from the
	// score instead of counting as wrong, and waits for a teacher.
	NeedsReview bool
	// AISuggestedCorrect is the verdict the assistant gave, nil when it could
	// not evaluate the file. Kept so the teacher sees who graded what.
	AISuggestedCorrect *bool
}

// evaluateAttachment grades the uploaded file with the assistant so the student
// gets an immediate result and can keep practising. What it could not read is
// left pending for a teacher.
//
// requireTeacher forces the pending path even when the assistant succeeded: on
// a level test the verdict decides promotion, so it is not left to the
// assistant alone.
func evaluateAttachment(
	ctx context.Context,
	app *appcontext.Context,
	cfg assistant.Config,
	ex domain.Exercise,
	gradeName string,
	attachmentURL, filename string,
	requireTeacher bool,
) attachmentOutcome {
	pending := attachmentOutcome{NeedsReview: true, Feedback: "Tu entrega quedó pendiente de revisión del docente."}

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
	// The submission body is user-controlled. The storage response is the
	// authoritative MIME for the object we actually fetched and evaluated.
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
	// The assistant received the file but could not make sense of it.
	if strings.Contains(strings.ToUpper(evaluation.Feedback), assistant.UnreadableFeedback) {
		return pending
	}

	verdict := evaluation.IsCorrect
	if requireTeacher {
		// The assistant's opinion is kept for the teacher, but it does not
		// count: promotion waits for a human.
		return attachmentOutcome{
			NeedsReview:        true,
			Feedback:           evaluation.Feedback,
			AISuggestedCorrect: &verdict,
		}
	}
	return attachmentOutcome{
		IsCorrect: verdict,
		Feedback:  evaluation.Feedback,
		// Graded, so the student is not blocked — but the teacher still sees it
		// in their queue and can change the grade.
		NeedsReview:        false,
		AISuggestedCorrect: &verdict,
	}
}

// attachmentKindAccepted reports whether the delivered file matches what the
// exercise allows. No configured list means any supported format.
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

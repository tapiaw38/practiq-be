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
	// Ungraded means nobody produced a verdict: the answer is excluded from the
	// score instead of counting as wrong. On a level test it also waits for a
	// teacher; a practice is never held up for one.
	Ungraded bool
	// AISuggestedCorrect is the verdict the assistant gave, nil when it could
	// not evaluate the file. Kept so the teacher sees who graded what.
	AISuggestedCorrect *bool
}

// evaluateAttachment grades the uploaded file with the assistant so the student
// gets an immediate result and can keep practising. What it could not read is
// left ungraded rather than counted as wrong.
//
// teacherGrades only picks the wording for an answer nobody could grade: a
// level test hands it to a teacher, a practice does not. It no longer forces
// the ungraded path — a file the assistant read and decided on counts on either
// sheet, and only what it could not resolve reaches a person.
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
	return attachmentOutcome{
		IsCorrect: verdict,
		Feedback:  evaluation.Feedback,
		// The assistant resolved it, so it counts — on a level test too. A
		// verdict it could reach is not a reason to hold the student's
		// promotion behind a person.
		Ungraded:           false,
		AISuggestedCorrect: &verdict,
	}
}

// ungradedAttachmentFeedback explains an answer that got no verdict. Only a
// level test hands it to a teacher; a practice says so plainly instead of
// promising a correction that is never going to arrive.
func ungradedAttachmentFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "Tu entrega quedó pendiente de revisión del docente."
	}
	return "No pudimos corregir esta entrega automáticamente, así que no cuenta en tu puntaje."
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

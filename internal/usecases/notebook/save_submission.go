package notebook

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/imagecompose"
	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

type (
	SaveSubmissionUsecase interface {
		Execute(ctx context.Context, input SaveSubmissionInput) error
	}

	SaveSubmissionInput struct {
		PageID     string
		StudentID  string
		CanvasData string
		AnswerText string
	}

	saveSubmissionUsecase struct{ contextFactory appcontext.Factory }
)

func NewSaveSubmissionUsecase(contextFactory appcontext.Factory) SaveSubmissionUsecase {
	return &saveSubmissionUsecase{contextFactory: contextFactory}
}

func (u *saveSubmissionUsecase) Execute(ctx context.Context, input SaveSubmissionInput) error {
	app := u.contextFactory()
	page, err := app.Repositories.Notebook.GetPage(ctx, input.PageID)
	if err != nil {
		return err
	}
	if page == nil {
		return fmt.Errorf("page not found")
	}
	notebook, err := app.Repositories.Notebook.Get(ctx, page.NotebookID)
	if err != nil {
		return err
	}
	if notebook == nil {
		return fmt.Errorf("notebook not found")
	}
	hasAccess, err := studentHasNotebookCourseAccess(ctx, app, input.StudentID, notebook.CourseID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("forbidden")
	}

	// Get course for grade context
	course, _ := app.Repositories.Course.Get(ctx, notebook.CourseID)
	gradeName := ""
	if course != nil {
		gradeName = course.GradeName
	}

	canvasForOCR := input.CanvasData
	submission := domain.NotebookSubmission{
		PageID:     input.PageID,
		StudentID:  input.StudentID,
		CanvasData: input.CanvasData,
		AnswerText: input.AnswerText,
	}

	profile, _ := app.Repositories.UserProfile.Get(ctx, input.StudentID)
	assistantCfg := assistant.Config{}
	if profile != nil {
		assistantCfg.BaseURL = profile.AssistantBaseURL
		assistantCfg.APIKey = profile.AssistantAPIKey
	}

	if page != nil && app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
		expectedAnswer := normalizeNotebookExpectedAnswer(page.ContentData)
		if expectedAnswer != "" {
			studentAnswer := strings.TrimSpace(input.AnswerText)
			if studentAnswer == "" && strings.TrimSpace(canvasForOCR) != "" {
				if resolved, err := resolveImageForOCR(ctx, app, canvasForOCR); err == nil {
					canvasForOCR = resolved
				} else {
					log.Printf("[image_storage] notebook submission resolve failed page_id=%s err=%v", input.PageID, err)
				}
				canvasForOCR = normalizeCanvasDataURI(canvasForOCR)
				if recognizedRaw, recognizeErr := app.Integrations.AssistantGateway.AnalyzeNotebookCanvas(ctx, assistantCfg, canvasForOCR, buildNotebookPromptContext(page)); recognizeErr == nil {
					recognizedText := strings.TrimSpace(recognizedRaw)
					submission.AIRecognizedText = recognizedText
					studentAnswer = recognizedText
				} else {
					submission.AIFeedback = "no se pudo analizar la imagen del cuaderno"
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
				}
			}

			if strings.EqualFold(studentAnswer, "UNREADABLE") {
				submission.AIFeedback = "respuesta no legible (UNREADABLE)"
				submission.AIReviewedAt = ptrTime(time.Now().UTC())
			}

			if studentAnswer != "" && studentAnswer != "UNREADABLE" {
				if evaluation, aiErr := evaluateNotebookSubmission(ctx, app, assistantCfg, page, canvasForOCR, expectedAnswer, studentAnswer, gradeName); aiErr == nil {
					submission.AIIsCorrect = &evaluation.IsCorrect
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
					if strings.TrimSpace(evaluation.Feedback) != "" {
						submission.AIFeedback = evaluation.Feedback
					} else if evaluation.IsCorrect {
						submission.AIFeedback = "respuesta evaluada como correcta"
					} else {
						submission.AIFeedback = "respuesta evaluada como incorrecta"
					}
				}
			} else {
				submission.AIFeedback = "no se pudo evaluar la respuesta"
				submission.AIReviewedAt = ptrTime(time.Now().UTC())
			}
		}
	}

	if isLikelyImageData(submission.CanvasData) && app.ImageStorage != nil {
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", input.StudentID, submission.CanvasData); err == nil {
			submission.CanvasData = uploaded
		} else {
			log.Printf("[image_storage] notebook submission upload failed page_id=%s student_id=%s err=%v", input.PageID, input.StudentID, err)
		}
	}

	return app.Repositories.Notebook.UpsertSubmission(ctx, submission)
}

// evaluateNotebookSubmission grades the student's work against the page it was
// set on.
//
// When the teacher's page is an image, there is no expected answer to compare
// text against — normalizeNotebookExpectedAnswer can only say "[imagen del
// docente]", so the grader was judging the transcription with no idea what the
// exercise asked. Stacking the two pages into one image and grading that gives
// it the consigna. Falls back to the text comparison whenever the pages cannot
// be composed, or when the page really is text.
func evaluateNotebookSubmission(
	ctx context.Context,
	app *appcontext.Context,
	cfg assistant.Config,
	page *domain.NotebookPage,
	studentCanvas, expectedAnswer, studentAnswer, gradeName string,
) (assistant.EvaluationResult, error) {
	pageContext := buildNotebookPromptContext(page)

	if composed := buildNotebookEvaluationImage(ctx, app, page, studentCanvas); composed != nil {
		return app.Integrations.AssistantGateway.EvaluateAttachment(ctx, cfg, assistant.AttachmentEvaluationInput{
			Kind:      assistant.AttachmentKindImage,
			Filename:  "cuaderno.png",
			Content:   composed,
			GradeName: gradeName,
			Question: pageContext +
				" La imagen adjunta tiene dos partes separadas por una linea horizontal:" +
				" arriba la pagina original del docente, que es la consigna;" +
				" abajo lo que resolvio el alumno." +
				" Evalua unicamente lo que escribio el alumno, comparandolo con la consigna de arriba.",
		})
	}

	return app.Integrations.AssistantGateway.EvaluatePracticeAnswer(
		ctx, cfg, pageContext, expectedAnswer, studentAnswer, gradeName,
	)
}

// buildNotebookEvaluationImage stacks the teacher's page above the student's
// work, or returns nil when either side is missing or cannot be decoded — in
// which case the caller grades on text alone rather than on half a picture.
func buildNotebookEvaluationImage(
	ctx context.Context,
	app *appcontext.Context,
	page *domain.NotebookPage,
	studentCanvas string,
) []byte {
	if page == nil {
		return nil
	}
	teacherPage := strings.TrimSpace(page.ContentData)
	if teacherPage == "" || !(isLikelyImageData(teacherPage) || isImageURL(teacherPage)) {
		return nil
	}

	teacherBytes := decodeNotebookImage(ctx, app, teacherPage)
	studentBytes := decodeNotebookImage(ctx, app, studentCanvas)
	if teacherBytes == nil || studentBytes == nil {
		return nil
	}

	composed, err := imagecompose.StackVertically(teacherBytes, studentBytes)
	if err != nil {
		log.Printf("[notebook] compose evaluation image failed page_id=%s err=%v", page.ID, err)
		return nil
	}
	return composed
}

func decodeNotebookImage(ctx context.Context, app *appcontext.Context, value string) []byte {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if resolved, err := resolveImageForOCR(ctx, app, value); err == nil {
		value = resolved
	}
	decoded, _, err := utils.DecodeDataURI(normalizeCanvasDataURI(value))
	if err != nil {
		return nil
	}
	return decoded
}

func buildNotebookPromptContext(page *domain.NotebookPage) string {
	if page == nil {
		return "Cuaderno"
	}
	return fmt.Sprintf("Cuaderno - Pagina %d. Titulo: %s. Instrucciones: %s", page.PageNumber, strings.TrimSpace(page.Title), strings.TrimSpace(page.Instructions))
}

func normalizeNotebookExpectedAnswer(contentData string) string {
	value := strings.TrimSpace(contentData)
	if value == "" {
		return ""
	}
	if isLikelyImageData(value) || isImageURL(value) {
		return "[imagen del docente]"
	}
	return value
}

// AddPage uploads a teacher's page image and stores the resulting URL, so
// ContentData is routinely an https link rather than the base64 isLikelyImageData
// looks for — and a URL is not base64-like, so it slipped through and reached the
// model verbatim as "respuesta correcta esperada: https://….png". Kept separate
// from isLikelyImageData because AddPage uses that one to decide what to upload,
// and an already-uploaded URL must not be uploaded again.
func isImageURL(value string) bool {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return false
	}
	if strings.ContainsAny(value, " \t\n") {
		return false
	}
	path := strings.ToLower(value)
	if idx := strings.IndexAny(path, "?#"); idx != -1 {
		path = path[:idx]
	}
	switch {
	case strings.HasSuffix(path, ".png"),
		strings.HasSuffix(path, ".jpg"),
		strings.HasSuffix(path, ".jpeg"),
		strings.HasSuffix(path, ".webp"),
		strings.HasSuffix(path, ".gif"):
		return true
	}
	return false
}

func isLikelyImageData(value string) bool {
	if strings.HasPrefix(value, "data:image/") {
		return true
	}
	compact := strings.ReplaceAll(strings.ReplaceAll(value, "\n", ""), "\r", "")
	if len(compact) < 128 {
		return false
	}
	if strings.HasPrefix(compact, "iVBORw0KGgo") || strings.HasPrefix(compact, "/9j/") || strings.HasPrefix(compact, "R0lGOD") {
		return true
	}
	if !isBase64Like(compact) {
		return false
	}
	return len(compact) > 512
}

func isBase64Like(value string) bool {
	for _, r := range value {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '+' || r == '/' || r == '=' {
			continue
		}
		return false
	}
	return true
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func normalizeCanvasDataURI(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(trimmed), "data:image/") {
		return trimmed
	}
	return "data:image/png;base64," + trimmed
}

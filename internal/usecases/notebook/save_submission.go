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
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	SaveSubmissionUsecase interface {
		Execute(ctx context.Context, pageID, studentID string, version int64, input SaveSubmissionInput) error
	}

	SaveSubmissionInput struct {
		CanvasData string `json:"canvas_data"`
		AnswerText string `json:"answer_text"`
	}

	saveSubmissionUsecase struct{ contextFactory appcontext.Factory }
)

func NewSaveSubmissionUsecase(contextFactory appcontext.Factory) SaveSubmissionUsecase {
	return &saveSubmissionUsecase{contextFactory: contextFactory}
}

func (u *saveSubmissionUsecase) Execute(ctx context.Context, pageID, studentID string, version int64, input SaveSubmissionInput) error {
	app := u.contextFactory()
	page, err := app.Repositories.Notebook.GetPage(ctx, pageID)
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
	hasAccess, err := studentHasNotebookCourseAccess(ctx, app, studentID, notebook.CourseID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("forbidden")
	}

	if appErr := school.EnsureCourseAcceptsWork(ctx, app, notebook.CourseID); appErr != nil {
		return fmt.Errorf("this course is closed: it no longer accepts submissions")
	}

	if appErr := school.EnsureStudentCanWork(ctx, app, studentID, notebook.CourseID); appErr != nil {
		return fmt.Errorf("tu docente pausó tu acceso: podés ver lo que ya hiciste, pero no entregar")
	}

	course, _ := app.Repositories.Course.Get(ctx, notebook.CourseID)
	gradeName := ""
	if course != nil {
		gradeName = course.GradeName
	}

	canvasForOCR := input.CanvasData
	submission := domain.NotebookSubmission{
		PageID:     pageID,
		StudentID:  studentID,
		CanvasData: input.CanvasData,
		AnswerText: input.AnswerText,
	}

	assistantCfg := assistantcfg.Resolve(ctx, app)

	statementReady := ensurePageStatement(ctx, app, notebook.TeacherID, page)

	hasStudentWork := strings.TrimSpace(input.AnswerText) != "" || strings.TrimSpace(input.CanvasData) != ""

	if statementReady && page != nil && app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
		expectedAnswer := normalizeNotebookExpectedAnswer(page.ContentData)
		if expectedAnswer != "" {
			studentAnswer := strings.TrimSpace(input.AnswerText)
			if studentAnswer == "" && strings.TrimSpace(canvasForOCR) != "" {
				if resolved, err := resolveImageForOCR(ctx, app, canvasForOCR); err == nil {
					canvasForOCR = resolved
				} else {
					log.Printf("[image_storage] notebook submission resolve failed page_id=%s err=%v", pageID, err)
				}
				canvasForOCR = normalizeCanvasDataURI(canvasForOCR)
				if recognizedRaw, recognizeErr := app.Integrations.AssistantGateway.AnalyzeNotebookCanvas(ctx, assistantCfg, canvasForOCR, buildNotebookPromptContext(page)); recognizeErr == nil {
					recognizedText := strings.TrimSpace(recognizedRaw)
					submission.AIRecognizedText = recognizedText
					studentAnswer = recognizedText
				} else {
					log.Printf("[notebook] canvas analysis failed page_id=%s err=%v", pageID, recognizeErr)
					submission.AIFeedback = "no se pudo analizar la imagen del cuaderno"
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
				}
			}

			if strings.EqualFold(studentAnswer, "UNREADABLE") {
				submission.AIFeedback = "respuesta no legible (UNREADABLE)"
				submission.AIReviewedAt = ptrTime(time.Now().UTC())
			} else if studentAnswer != "" {
				if evaluation, aiErr := evaluateNotebookSubmission(ctx, app, assistantCfg, page, expectedAnswer, studentAnswer, gradeName); aiErr == nil {
					submission.AIIsCorrect = &evaluation.IsCorrect
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
					if strings.TrimSpace(evaluation.Feedback) != "" {
						submission.AIFeedback = evaluation.Feedback
					} else if evaluation.IsCorrect {
						submission.AIFeedback = "respuesta evaluada como correcta"
					} else {
						submission.AIFeedback = "respuesta evaluada como incorrecta"
					}
				} else {
					log.Printf("[notebook] evaluation failed page_id=%s err=%v", pageID, aiErr)
					submission.AIFeedback = "no se pudo evaluar la respuesta"
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
				}
			} else {
				submission.AIFeedback = "no se encontro respuesta para evaluar"
				submission.AIReviewedAt = ptrTime(time.Now().UTC())
			}
		}
	}

	submission.NeedsTeacherReview = submissionNeedsTeacherReview(hasStudentWork, submission.AIIsCorrect)
	submission.Version = version

	if isLikelyImageData(submission.CanvasData) && app.ImageStorage != nil {
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", studentID, submission.CanvasData); err == nil {
			submission.CanvasData = uploaded
		} else {
			log.Printf("[image_storage] notebook submission upload failed page_id=%s student_id=%s err=%v", pageID, studentID, err)
		}
	}

	return app.Repositories.Notebook.UpsertSubmission(ctx, submission)
}

func evaluateNotebookSubmission(
	ctx context.Context,
	app *appcontext.Context,
	cfg assistant.Config,
	page *domain.NotebookPage,
	expectedAnswer, studentAnswer, gradeName string,
) (assistant.EvaluationResult, error) {
	return app.Integrations.AssistantGateway.EvaluatePracticeAnswer(
		ctx, cfg, buildNotebookPromptContext(page), expectedAnswer, studentAnswer, gradeName,
	)
}

func submissionNeedsTeacherReview(hasStudentWork bool, aiIsCorrect *bool) bool {
	return hasStudentWork && aiIsCorrect == nil
}

func buildNotebookPromptContext(page *domain.NotebookPage) string {
	if page == nil {
		return "Cuaderno"
	}

	context := fmt.Sprintf(
		"Cuaderno - Pagina %d. Titulo: %s. Instrucciones: %s",
		page.PageNumber,
		strings.TrimSpace(page.Title),
		strings.TrimSpace(page.Instructions),
	)

	if statement := strings.TrimSpace(page.StatementText); statement != "" {
		context += " Consigna de la pagina: " + statement
	}
	return context
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

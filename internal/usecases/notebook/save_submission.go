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
				if recognizedRaw, recognizeErr := app.Integrations.AssistantGateway.AnalyzeCanvas(ctx, assistantCfg, canvasForOCR, expectedAnswer); recognizeErr == nil {
					recognizedText := strings.TrimSpace(recognizedRaw)
					submission.AIRecognizedText = recognizedText
					studentAnswer = recognizedText
				} else {
					submission.AIFeedback = "Gillie: no se pudo analizar la imagen del cuaderno"
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
				}
			}

			if strings.EqualFold(studentAnswer, "UNREADABLE") {
				submission.AIFeedback = "Gillie: respuesta no legible (UNREADABLE)"
				submission.AIReviewedAt = ptrTime(time.Now().UTC())
			}

			if studentAnswer != "" && studentAnswer != "UNREADABLE" {
				if evaluation, aiErr := app.Integrations.AssistantGateway.EvaluatePracticeAnswer(ctx, assistantCfg, buildNotebookPromptContext(page), expectedAnswer, studentAnswer); aiErr == nil {
					submission.AIIsCorrect = &evaluation.IsCorrect
					submission.AIReviewedAt = ptrTime(time.Now().UTC())
					if strings.TrimSpace(evaluation.Feedback) != "" {
						submission.AIFeedback = evaluation.Feedback
					} else if evaluation.IsCorrect {
						submission.AIFeedback = "Gillie: respuesta evaluada como correcta"
					} else {
						submission.AIFeedback = "Gillie: respuesta evaluada como incorrecta"
					}
				}
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
	if isLikelyImageData(value) {
		return "[imagen del docente]"
	}
	return value
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

package notebook

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/usecases/datauri"
)

type (
	ReviewSubmissionUsecase interface {
		Execute(ctx context.Context, submissionID string, teacherID string) (*ReviewSubmissionOutput, error)
	}

	ReviewSubmissionOutput struct {
		Data NotebookSubmissionFullData `json:"data"`
	}

	reviewSubmissionUsecase struct{ contextFactory appcontext.Factory }
)

func NewReviewSubmissionUsecase(contextFactory appcontext.Factory) ReviewSubmissionUsecase {
	return &reviewSubmissionUsecase{contextFactory: contextFactory}
}

func (u *reviewSubmissionUsecase) Execute(ctx context.Context, submissionID string, teacherID string) (*ReviewSubmissionOutput, error) {
	app := u.contextFactory()

	submission, err := app.Repositories.Notebook.GetFullSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if submission == nil || (teacherID != "" && submission.TeacherID != teacherID) {
		return nil, fmt.Errorf("submission not found")
	}

	page, err := app.Repositories.Notebook.GetPage(ctx, submission.PageID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, fmt.Errorf("page not found")
	}

	profile, _ := app.Repositories.UserProfile.Get(ctx, submission.StudentID)
	assistantCfg := assistant.Config{}
	if profile != nil {
		assistantCfg.BaseURL = profile.AssistantBaseURL
		assistantCfg.APIKey = profile.AssistantAPIKey
	}

	if app.Integrations.AssistantGateway == nil || !app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
		return nil, fmt.Errorf("assistant service not configured")
	}

	// Get course for grade context
	course, _ := app.Repositories.Course.Get(ctx, submission.CourseID)
	gradeName := ""
	if course != nil {
		gradeName = course.GradeName
	}

	expectedAnswer := normalizeNotebookExpectedAnswer(page.ContentData)
	if expectedAnswer == "" {
		expectedAnswer = "[sin respuesta esperada]"
	}

	var recognizedText string
	var isCorrect *bool
	var feedback string
	ensurePageStatement(ctx, app, assistantCfg, page)
	needsTeacherReview := statementNeedsTeacherReview(page)

	studentAnswer := strings.TrimSpace(submission.AnswerText)

	// If there's canvas data but no text answer, try to recognize the handwriting
	if studentAnswer == "" && strings.TrimSpace(submission.CanvasData) != "" {
		canvasData, resolveErr := resolveImageForOCR(ctx, app, submission.CanvasData)
		if resolveErr != nil {
			log.Printf("[image_storage] notebook review resolve failed submission_id=%s err=%v", submissionID, resolveErr)
			canvasData = submission.CanvasData
		}
		canvasData = normalizeCanvasDataURI(canvasData)
		recognized, recognizeErr := app.Integrations.AssistantGateway.AnalyzeNotebookCanvas(ctx, assistantCfg, canvasData, buildNotebookPromptContext(page))
		if recognizeErr != nil {
			log.Printf("[notebook] canvas analysis failed submission_id=%s err=%v", submissionID, recognizeErr)
			needsTeacherReview = true
			feedback = "no se pudo analizar la imagen del cuaderno"
		} else {
			recognizedText = strings.TrimSpace(recognized)
			studentAnswer = recognizedText
		}
	}

	// If the recognized answer is unreadable, report that
	if strings.EqualFold(studentAnswer, "UNREADABLE") {
		needsTeacherReview = true
		feedback = "respuesta no legible (UNREADABLE)"
		isCorrect = nil
	} else if studentAnswer != "" {
		// Evaluate the answer with AI, against the teacher's page when it is an
		// image (see evaluateNotebookSubmission).
		evaluation, aiErr := evaluateNotebookSubmission(ctx, app, assistantCfg, page, expectedAnswer, studentAnswer, gradeName)
		if aiErr != nil {
			log.Printf("[notebook] evaluation failed submission_id=%s err=%v", submissionID, aiErr)
			needsTeacherReview = true
			feedback = "no se pudo evaluar la respuesta"
		} else {
			isCorrect = &evaluation.IsCorrect
			if strings.TrimSpace(evaluation.Feedback) != "" {
				feedback = evaluation.Feedback
			} else if evaluation.IsCorrect {
				feedback = "respuesta evaluada como correcta"
			} else {
				feedback = "respuesta evaluada como incorrecta"
			}
		}
	} else {
		needsTeacherReview = false
		feedback = "no se encontro respuesta para evaluar"
	}

	// Update the submission with the AI review results
	if err := app.Repositories.Notebook.UpdateSubmissionAIReview(ctx, submissionID, recognizedText, isCorrect, feedback, needsTeacherReview); err != nil {
		return nil, err
	}

	updated, err := app.Repositories.Notebook.GetFullSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, fmt.Errorf("submission not found")
	}
	updated.CanvasData = datauri.Resolve(ctx, app, updated.CanvasData)
	return &ReviewSubmissionOutput{Data: toFullSubmissionData(*updated)}, nil
}

func resolveImageForOCR(ctx context.Context, app *appcontext.Context, value string) (string, error) {
	if app.ImageStorage == nil {
		return value, nil
	}
	return app.ImageStorage.ResolveDataURI(ctx, value)
}

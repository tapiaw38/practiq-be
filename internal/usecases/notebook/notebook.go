package notebook

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/assistant"
)

// ── Create ──────────────────────────────────────────────────────────

type CreateUsecase interface {
	Execute(ctx context.Context, teacherID string, input CreateInput) (*NotebookOutput, error)
}

type CreateInput struct {
	CourseID    string
	Title       string
	Description string
	Level       int
}

type createUsecase struct{ factory appcontext.Factory }

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, teacherID string, input CreateInput) (*NotebookOutput, error) {
	app := u.factory()
	id, err := app.Repositories.Notebook.Create(ctx, domain.Notebook{
		CourseID:    input.CourseID,
		TeacherID:   teacherID,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
	})
	if err != nil {
		return nil, err
	}
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toOutput(nb), nil
}

// ── List ─────────────────────────────────────────────────────────────

type ListUsecase interface {
	Execute(ctx context.Context, courseID string) ([]NotebookOutput, error)
}

type listUsecase struct{ factory appcontext.Factory }

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) ([]NotebookOutput, error) {
	app := u.factory()
	notebooks, err := app.Repositories.Notebook.List(ctx, courseID)
	if err != nil {
		return nil, err
	}
	out := make([]NotebookOutput, 0, len(notebooks))
	for _, nb := range notebooks {
		nb := nb
		out = append(out, *toOutput(&nb))
	}
	return out, nil
}

// ── Get ──────────────────────────────────────────────────────────────

type GetUsecase interface {
	Execute(ctx context.Context, id, studentID string) (*NotebookOutput, error)
}

type getUsecase struct{ factory appcontext.Factory }

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, id, studentID string) (*NotebookOutput, error) {
	app := u.factory()
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if nb == nil {
		return nil, nil
	}

	// Attach student submissions to each page
	if studentID != "" {
		for i := range nb.Pages {
			sub, _ := app.Repositories.Notebook.GetSubmission(ctx, nb.Pages[i].ID, studentID)
			nb.Pages[i].Submission = sub
		}
	}

	return toOutput(nb), nil
}

// ── Add Page ─────────────────────────────────────────────────────────

type AddPageUsecase interface {
	Execute(ctx context.Context, input AddPageInput) (*PageOutput, error)
}

type AddPageInput struct {
	NotebookID   string
	PageNumber   int
	Title        string
	ContentType  string
	ContentData  string
	Instructions string
}

type addPageUsecase struct{ factory appcontext.Factory }

func NewAddPageUsecase(factory appcontext.Factory) AddPageUsecase {
	return &addPageUsecase{factory: factory}
}

func (u *addPageUsecase) Execute(ctx context.Context, input AddPageInput) (*PageOutput, error) {
	app := u.factory()
	id, err := app.Repositories.Notebook.CreatePage(ctx, domain.NotebookPage{
		NotebookID:   input.NotebookID,
		PageNumber:   input.PageNumber,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  input.ContentData,
		Instructions: input.Instructions,
	})
	if err != nil {
		return nil, err
	}
	return &PageOutput{ID: id, NotebookID: input.NotebookID, PageNumber: input.PageNumber,
		Title: input.Title, ContentType: input.ContentType, ContentData: input.ContentData,
		Instructions: input.Instructions}, nil
}

// ── Update Page ───────────────────────────────────────────────────────

type UpdatePageUsecase interface {
	Execute(ctx context.Context, input UpdatePageInput) error
}

type UpdatePageInput struct {
	PageID       string
	Title        string
	ContentType  string
	ContentData  string
	Instructions string
}

type updatePageUsecase struct{ factory appcontext.Factory }

func NewUpdatePageUsecase(factory appcontext.Factory) UpdatePageUsecase {
	return &updatePageUsecase{factory: factory}
}

func (u *updatePageUsecase) Execute(ctx context.Context, input UpdatePageInput) error {
	app := u.factory()
	return app.Repositories.Notebook.UpdatePage(ctx, domain.NotebookPage{
		ID:           input.PageID,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  input.ContentData,
		Instructions: input.Instructions,
	})
}

// ── Save Submission ───────────────────────────────────────────────────

type SaveSubmissionUsecase interface {
	Execute(ctx context.Context, input SaveSubmissionInput) error
}

type SaveSubmissionInput struct {
	PageID     string
	StudentID  string
	CanvasData string
	AnswerText string
}

type saveSubmissionUsecase struct{ factory appcontext.Factory }

func NewSaveSubmissionUsecase(factory appcontext.Factory) SaveSubmissionUsecase {
	return &saveSubmissionUsecase{factory: factory}
}

func (u *saveSubmissionUsecase) Execute(ctx context.Context, input SaveSubmissionInput) error {
	app := u.factory()
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

	page, _ := app.Repositories.Notebook.GetPage(ctx, input.PageID)
	if page != nil && app.AssistantService != nil && app.AssistantService.IsConfigured(assistantCfg) {
		expectedAnswer := normalizeNotebookExpectedAnswer(page.ContentData)
		if expectedAnswer != "" {
			studentAnswer := strings.TrimSpace(input.AnswerText)
			if studentAnswer == "" && strings.TrimSpace(input.CanvasData) != "" {
				if recognizedRaw, recognizeErr := app.AssistantService.AnalyzeCanvas(ctx, assistantCfg, input.CanvasData, expectedAnswer); recognizeErr == nil {
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
				if evaluation, aiErr := app.AssistantService.EvaluatePracticeAnswer(ctx, assistantCfg, buildNotebookPromptContext(page), expectedAnswer, studentAnswer); aiErr == nil {
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

// ── Update Notebook ───────────────────────────────────────────────────

type UpdateUsecase interface {
	Execute(ctx context.Context, id string, input UpdateInput) (*NotebookOutput, error)
}

type UpdateInput struct {
	Title       string
	Description string
}

type updateUsecase struct{ factory appcontext.Factory }

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*NotebookOutput, error) {
	app := u.factory()
	if err := app.Repositories.Notebook.Update(ctx, id, domain.Notebook{
		Title:       input.Title,
		Description: input.Description,
	}); err != nil {
		return nil, err
	}
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toOutput(nb), nil
}

// ── Delete Notebook ───────────────────────────────────────────────────

type DeleteUsecase interface {
	Execute(ctx context.Context, id string) error
}

type deleteUsecase struct{ factory appcontext.Factory }

func NewDeleteUsecase(factory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{factory: factory}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) error {
	app := u.factory()
	return app.Repositories.Notebook.Delete(ctx, id)
}

// ── Output types ──────────────────────────────────────────────────────

type SubmissionOutput struct {
	ID               string `json:"id"`
	CanvasData       string `json:"canvas_data"`
	AnswerText       string `json:"answer_text"`
	AIRecognizedText string `json:"ai_recognized_text,omitempty"`
	AIIsCorrect      *bool  `json:"ai_is_correct,omitempty"`
	AIFeedback       string `json:"ai_feedback,omitempty"`
	AIReviewedAt     string `json:"ai_reviewed_at,omitempty"`
}

type PageOutput struct {
	ID           string            `json:"id"`
	NotebookID   string            `json:"notebook_id"`
	PageNumber   int               `json:"page_number"`
	Title        string            `json:"title"`
	ContentType  string            `json:"content_type"`
	ContentData  string            `json:"content_data"`
	Instructions string            `json:"instructions"`
	Submission   *SubmissionOutput `json:"submission,omitempty"`
}

type NotebookOutput struct {
	Data NotebookData `json:"data"`
}

type NotebookData struct {
	ID          string       `json:"id"`
	CourseID    string       `json:"course_id"`
	TeacherID   string       `json:"teacher_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Level       int          `json:"level"`
	Pages       []PageOutput `json:"pages"`
	CreatedAt   string       `json:"created_at"`
}

func toOutput(nb *domain.Notebook) *NotebookOutput {
	pages := make([]PageOutput, 0, len(nb.Pages))
	for _, p := range nb.Pages {
		po := PageOutput{
			ID:           p.ID,
			NotebookID:   p.NotebookID,
			PageNumber:   p.PageNumber,
			Title:        p.Title,
			ContentType:  p.ContentType,
			ContentData:  p.ContentData,
			Instructions: p.Instructions,
		}
		if p.Submission != nil {
			aiReviewedAt := ""
			if p.Submission.AIReviewedAt != nil {
				aiReviewedAt = p.Submission.AIReviewedAt.Format("2006-01-02T15:04:05Z")
			}
			po.Submission = &SubmissionOutput{
				ID:               p.Submission.ID,
				CanvasData:       p.Submission.CanvasData,
				AnswerText:       p.Submission.AnswerText,
				AIRecognizedText: p.Submission.AIRecognizedText,
				AIIsCorrect:      p.Submission.AIIsCorrect,
				AIFeedback:       p.Submission.AIFeedback,
				AIReviewedAt:     aiReviewedAt,
			}
		}
		pages = append(pages, po)
	}
	return &NotebookOutput{Data: NotebookData{
		ID:          nb.ID,
		CourseID:    nb.CourseID,
		TeacherID:   nb.TeacherID,
		Title:       nb.Title,
		Description: nb.Description,
		Level:       nb.Level,
		Pages:       pages,
		CreatedAt:   nb.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}}
}

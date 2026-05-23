package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
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
	return app.Repositories.Notebook.UpsertSubmission(ctx, domain.NotebookSubmission{
		PageID:     input.PageID,
		StudentID:  input.StudentID,
		CanvasData: input.CanvasData,
		AnswerText: input.AnswerText,
	})
}

// ── Output types ──────────────────────────────────────────────────────

type SubmissionOutput struct {
	ID         string `json:"id"`
	CanvasData string `json:"canvas_data"`
	AnswerText string `json:"answer_text"`
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
			po.Submission = &SubmissionOutput{
				ID:         p.Submission.ID,
				CanvasData: p.Submission.CanvasData,
				AnswerText: p.Submission.AnswerText,
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

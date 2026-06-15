package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	CreateUsecase interface {
		Execute(ctx context.Context, teacherID string, input CreateInput) (*CreateOutput, error)
	}

	CreateInput struct {
		CourseID    string
		Title       string
		Description string
		Level       int
	}

	CreateOutput struct {
		Data NotebookData `json:"data"`
	}

	createUsecase struct{ contextFactory appcontext.Factory }
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, teacherID string, input CreateInput) (*CreateOutput, error) {
	app := u.contextFactory()
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
	return &CreateOutput{Data: toNotebookData(nb)}, nil
}

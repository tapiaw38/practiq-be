package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError)
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

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	course, err := app.Repositories.Course.Get(ctx, input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}
	if !isAdmin && course.TeacherID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	id, err := app.Repositories.Notebook.Create(ctx, domain.Notebook{
		CourseID:    input.CourseID,
		TeacherID:   requesterID,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	return &CreateOutput{Data: toNotebookData(nb)}, nil
}

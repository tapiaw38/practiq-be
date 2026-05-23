package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*GradeOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	Name        string
	Description string
	CreatedBy   string
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*GradeOutput, apperrors.ApplicationError) {
	app := u.factory()

	id, err := app.Repositories.Grade.Create(ctx, domain.Grade{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeCreateError, err)
	}

	grade, err := app.Repositories.Grade.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	return &GradeOutput{Data: toGradeData(*grade)}, nil
}

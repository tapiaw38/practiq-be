package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*SubjectOutput, apperrors.ApplicationError)
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

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*SubjectOutput, apperrors.ApplicationError) {
	app := u.factory()
	id, err := app.Repositories.Subject.Create(ctx, domain.Subject{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectCreateError, err)
	}
	subject, err := app.Repositories.Subject.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectGetError, err)
	}
	if subject == nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectNotFoundError, nil)
	}
	return &SubjectOutput{Data: toSubjectData(*subject)}, nil
}

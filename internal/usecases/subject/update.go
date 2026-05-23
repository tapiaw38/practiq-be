package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*SubjectOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	Name        string
	Description string
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*SubjectOutput, apperrors.ApplicationError) {
	app := u.factory()

	if err := app.Repositories.Subject.Update(ctx, id, domain.Subject{
		Name:        input.Name,
		Description: input.Description,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectUpdateError, err)
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

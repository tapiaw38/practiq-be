package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(context.Context, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Name        string
		Description string
	}

	UpdateOutput struct {
		Data SubjectData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &UpdateOutput{Data: toSubjectData(*subject)}, nil
}

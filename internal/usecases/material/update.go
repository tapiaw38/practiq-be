package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*MaterialOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	Title         string
	ExtractedText string
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*MaterialOutput, apperrors.ApplicationError) {
	app := u.factory()

	if err := app.Repositories.Material.Update(ctx, id, domain.Material{
		Title:         input.Title,
		ExtractedText: input.ExtractedText,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialUpdateError, err)
	}

	m, err := app.Repositories.Material.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if m == nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialNotFoundError, nil)
	}

	return &MaterialOutput{Data: toMaterialData(*m)}, nil
}

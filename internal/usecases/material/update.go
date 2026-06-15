package material

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
		Title         string
		ExtractedText string
	}

	UpdateOutput struct {
		Data MaterialData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &UpdateOutput{Data: toMaterialData(*m)}, nil
}

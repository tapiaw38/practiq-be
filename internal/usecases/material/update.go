package material

import (
	"context"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Title         string
		ExtractedText string
		FileURL       string
	}

	UpdateOutput struct {
		Data MaterialData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify material exists and check ownership
	material, err := app.Repositories.Material.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if material == nil {
		return nil, apperrors.NewNotFoundError("material not found")
	}

	if !isAdmin && material.TeacherID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	// Same rule as create, except an unchanged URL always passes: the file may
	// predate this check, or belong to another teacher of the same course, and
	// re-saving a title must not fail because of it.
	if input.FileURL != "" && input.FileURL != material.FileURL &&
		(app.ImageStorage == nil ||
			!app.ImageStorage.OwnsFileURL(input.FileURL, materialsFolder, requesterID)) {
		return nil, apperrors.NewApplicationError(mappings.MaterialUpdateError,
			errors.New("the file does not belong to this teacher"))
	}

	if err := app.Repositories.Material.Update(ctx, id, domain.Material{
		Title:         input.Title,
		ExtractedText: input.ExtractedText,
		FileURL:       input.FileURL,
		Status:        materialStatus(input.FileURL),
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

	return &UpdateOutput{Data: withViewURL(app, toMaterialData(*m))}, nil
}

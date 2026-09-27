package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, materialID string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data MaterialData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, materialID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	material, err := app.Repositories.Material.Get(ctx, materialID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if material == nil {
		return nil, apperrors.NewNotFoundError("material not found")
	}

	if appErr := requesterCanReadCourse(ctx, app, requesterID, isSuperAdmin, material.CourseID); appErr != nil {
		return nil, appErr
	}

	return &GetOutput{Data: withViewURL(app, toMaterialData(*material))}, nil
}

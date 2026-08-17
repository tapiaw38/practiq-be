package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	// GetUsecase reads one material, extracted text included.
	//
	// It is the other half of the listing: that one carries a preview so a
	// course's materials do not ship every document's full text, and this one
	// serves the rest when a reader actually opens one.
	GetUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, materialID string) (*GetOutput, apperrors.ApplicationError)
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

func (u *getUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, materialID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	material, err := app.Repositories.Material.Get(ctx, materialID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if material == nil {
		return nil, apperrors.NewNotFoundError("material not found")
	}

	// Reading a material is reading its course. Checking after the lookup is
	// what lets the course be known at all; the material id alone says nothing.
	if appErr := requesterCanReadCourse(ctx, app, requesterID, isAdmin, material.CourseID); appErr != nil {
		return nil, appErr
	}

	return &GetOutput{Data: withViewURL(app, toMaterialData(*material))}, nil
}

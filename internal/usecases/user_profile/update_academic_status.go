package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateAcademicStatusUsecase interface {
		Execute(context.Context, string, string) (*UpdateAcademicStatusOutput, apperrors.ApplicationError)
	}

	updateAcademicStatusUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateAcademicStatusOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateAcademicStatusUsecase(contextFactory appcontext.Factory) UpdateAcademicStatusUsecase {
	return &updateAcademicStatusUsecase{contextFactory: contextFactory}
}

func (u *updateAcademicStatusUsecase) Execute(ctx context.Context, id, status string) (*UpdateAcademicStatusOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if status != "active" && status != "blocked" {
		return nil, apperrors.NewBadRequestError("academic_status must be active or blocked")
	}

	if err := app.Repositories.UserProfile.UpdateAcademicStatus(ctx, id, status); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileUpdateError, err)
	}

	profile, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewApplicationError(mappings.NotFoundError, nil)
	}

	return &UpdateAcademicStatusOutput{Data: toProfileData(*profile)}, nil
}

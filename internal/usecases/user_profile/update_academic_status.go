package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateAcademicStatusUsecase interface {
	Execute(context.Context, string, string) (*ProfileOutput, apperrors.ApplicationError)
}

type updateAcademicStatusUsecase struct {
	factory appcontext.Factory
}

func NewUpdateAcademicStatusUsecase(factory appcontext.Factory) UpdateAcademicStatusUsecase {
	return &updateAcademicStatusUsecase{factory: factory}
}

func (u *updateAcademicStatusUsecase) Execute(ctx context.Context, id, status string) (*ProfileOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	return &ProfileOutput{Data: toProfileData(*profile)}, nil
}

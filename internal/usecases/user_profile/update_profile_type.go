package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	UpdateProfileTypeUsecase interface {
		Execute(context.Context, string, string, string) (*UpdateProfileTypeOutput, apperrors.ApplicationError)
	}

	updateProfileTypeUsecase struct{ contextFactory appcontext.Factory }

	UpdateProfileTypeOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateProfileTypeUsecase(contextFactory appcontext.Factory) UpdateProfileTypeUsecase {
	return &updateProfileTypeUsecase{contextFactory: contextFactory}
}

func (u *updateProfileTypeUsecase) Execute(ctx context.Context, id, profileType, bearerToken string) (*UpdateProfileTypeOutput, apperrors.ApplicationError) {
	if profileType != "teacher" && profileType != "student" {
		return nil, apperrors.NewBadRequestError("profile_type must be teacher or student")
	}
	app := u.contextFactory()
	profile, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewApplicationError(mappings.NotFoundError, nil)
	}
	if err := app.Repositories.UserProfile.UpdateProfileType(ctx, id, profileType); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileUpdateError, err)
	}
	profile, err = app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{id})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	info := names[id]
	return &UpdateProfileTypeOutput{Data: toProfileData(*profile, identity.FullName(info, id), info.Email)}, nil
}

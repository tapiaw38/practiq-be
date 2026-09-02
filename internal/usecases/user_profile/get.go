package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, id, bearerToken string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, id, bearerToken string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	p, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if p == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{id})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	info := names[id]

	return &GetOutput{Data: toProfileData(*p, identity.FullName(info, id), info.Email)}, nil
}

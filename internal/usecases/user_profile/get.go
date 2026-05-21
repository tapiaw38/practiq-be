package userprofile

import (
	"context"

	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type GetUsecase interface {
	Execute(context.Context, string) (*ProfileOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	factory appcontext.Factory
}

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*ProfileOutput, apperrors.ApplicationError) {
	app := u.factory()

	p, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if p == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	return &ProfileOutput{Data: toProfileData(*p)}, nil
}

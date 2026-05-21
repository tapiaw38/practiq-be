package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type SyncUsecase interface {
	Execute(context.Context, SyncInput) (*ProfileOutput, apperrors.ApplicationError)
}

type syncUsecase struct {
	factory appcontext.Factory
}

type SyncInput struct {
	ID          string
	Name        string
	Email       string
	ProfileType string
}

func NewSyncUsecase(factory appcontext.Factory) SyncUsecase {
	return &syncUsecase{factory: factory}
}

func (u *syncUsecase) Execute(ctx context.Context, input SyncInput) (*ProfileOutput, apperrors.ApplicationError) {
	app := u.factory()

	profileType := input.ProfileType
	if profileType == "" {
		profileType = "student"
	}
	if profileType != "teacher" && profileType != "student" {
		return nil, apperrors.NewBadRequestError("profile_type must be teacher or student")
	}

	p := domain.UserProfile{
		ID:          input.ID,
		Name:        input.Name,
		Email:       input.Email,
		ProfileType: profileType,
	}

	if err := app.Repositories.UserProfile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	return &ProfileOutput{Data: toProfileData(*updated)}, nil
}

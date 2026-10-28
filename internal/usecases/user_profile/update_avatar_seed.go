package userprofile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

const maxAvatarSeedLength = 64

type (
	UpdateAvatarSeedUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateAvatarSeedInput) (*UpdateAvatarSeedOutput, apperrors.ApplicationError)
	}

	updateAvatarSeedUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateAvatarSeedInput struct {
		AvatarSeed string `json:"avatar_seed"`
	}

	UpdateAvatarSeedOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateAvatarSeedUsecase(contextFactory appcontext.Factory) UpdateAvatarSeedUsecase {
	return &updateAvatarSeedUsecase{contextFactory: contextFactory}
}

func (u *updateAvatarSeedUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateAvatarSeedInput) (*UpdateAvatarSeedOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}

	seed := strings.TrimSpace(in.AvatarSeed)
	if !validAvatarSeed(seed) {
		return nil, apperrors.NewBadRequestError("avatar_seed must be up to 64 letters, digits, hyphens or underscores")
	}

	if err := app.Repositories.UserProfile.UpdateAvatarSeed(ctx, id, seed); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{id})
	if appErr != nil {
		return nil, appErr
	}
	info := names[id]

	return &UpdateAvatarSeedOutput{Data: toProfileData(*updated, identity.FullName(info, id), info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

func validAvatarSeed(seed string) bool {
	if len(seed) > maxAvatarSeedLength {
		return false
	}
	for _, r := range seed {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_':
		default:
			return false
		}
	}
	return true
}

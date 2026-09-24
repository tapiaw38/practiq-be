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
		Execute(context.Context, string, bool, UpdateAvatarSeedInput) (*UpdateAvatarSeedOutput, apperrors.ApplicationError)
	}

	updateAvatarSeedUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateAvatarSeedInput struct {
		ID          string
		AvatarSeed  string `json:"avatar_seed"`
		BearerToken string
	}

	UpdateAvatarSeedOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateAvatarSeedUsecase(contextFactory appcontext.Factory) UpdateAvatarSeedUsecase {
	return &updateAvatarSeedUsecase{contextFactory: contextFactory}
}

func (u *updateAvatarSeedUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input UpdateAvatarSeedInput) (*UpdateAvatarSeedOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Same rule as the other profile writes: your own profile, or one in a
	// school you administer.
	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, input.ID); appErr != nil {
		return nil, appErr
	}

	seed := strings.TrimSpace(input.AvatarSeed)
	if !validAvatarSeed(seed) {
		return nil, apperrors.NewBadRequestError("avatar_seed must be up to 64 letters, digits, hyphens or underscores")
	}

	if err := app.Repositories.UserProfile.UpdateAvatarSeed(ctx, input.ID, seed); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, []string{input.ID})
	if appErr != nil {
		return nil, appErr
	}
	info := names[input.ID]

	return &UpdateAvatarSeedOutput{Data: toProfileData(*updated, identity.FullName(info, input.ID), info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

// validAvatarSeed keeps the seed to characters that are safe to drop straight
// into a URL or an SVG id. Empty is allowed: it clears the avatar back to the
// default. The point is that the stored value can never be markup, a URL or
// anything the client would resolve — only an opaque token to draw from.
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

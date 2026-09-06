package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	SyncUsecase interface {
		Execute(context.Context, SyncInput) (*SyncOutput, apperrors.ApplicationError)
	}

	syncUsecase struct {
		contextFactory appcontext.Factory
	}

	SyncInput struct {
		ID               string
		ProfileType      string
		Timezone         string
		AssistantBaseURL string
		AssistantAPIKey  string
		// BearerToken is the caller's own "Bearer <jwt>" header, forwarded to
		// auth-api-be to resolve the caller's own display name — never
		// trusted from the request body, since that would let the client
		// claim any name it likes.
		BearerToken string
	}

	SyncOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewSyncUsecase(contextFactory appcontext.Factory) SyncUsecase {
	return &syncUsecase{contextFactory: contextFactory}
}

func (u *syncUsecase) Execute(ctx context.Context, input SyncInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	profileType := input.ProfileType
	if profileType == "" {
		profileType = "student"
		if existing, err := app.Repositories.UserProfile.Get(ctx, input.ID); err == nil && existing != nil && existing.ProfileType != "" {
			profileType = existing.ProfileType
		}
	}
	if profileType != "teacher" && profileType != "student" {
		return nil, apperrors.NewBadRequestError("profile_type must be teacher or student")
	}

	p := domain.UserProfile{
		ID:               input.ID,
		ProfileType:      profileType,
		Timezone:         input.Timezone,
		AssistantBaseURL: input.AssistantBaseURL,
		AssistantAPIKey:  input.AssistantAPIKey,
	}

	if err := app.Repositories.UserProfile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, []string{input.ID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	info := names[input.ID]

	return &SyncOutput{Data: toProfileData(*updated, identity.FullName(info, input.ID), info.Email)}, nil
}

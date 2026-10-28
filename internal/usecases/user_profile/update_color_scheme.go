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

type (
	UpdateColorSchemeUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateColorSchemeInput) (*UpdateColorSchemeOutput, apperrors.ApplicationError)
	}

	updateColorSchemeUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateColorSchemeInput struct {
		ColorScheme string `json:"color_scheme"`
	}

	UpdateColorSchemeOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateColorSchemeUsecase(contextFactory appcontext.Factory) UpdateColorSchemeUsecase {
	return &updateColorSchemeUsecase{contextFactory: contextFactory}
}

func (u *updateColorSchemeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateColorSchemeInput) (*UpdateColorSchemeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}

	scheme := strings.TrimSpace(in.ColorScheme)
	if scheme == "" {
		scheme = "light"
	}
	if scheme != "light" && scheme != "dark" {
		return nil, apperrors.NewBadRequestError("color_scheme must be light or dark")
	}

	if err := app.Repositories.UserProfile.UpdateColorScheme(ctx, id, scheme); err != nil {
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

	return &UpdateColorSchemeOutput{Data: toProfileData(*updated, identity.FullName(info, id), info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

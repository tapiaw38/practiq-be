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
	UpdateUIThemeUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateUIThemeInput) (*UpdateUIThemeOutput, apperrors.ApplicationError)
	}

	updateUIThemeUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateUIThemeInput struct {
		UITheme string `json:"ui_theme"`
	}

	UpdateUIThemeOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateUIThemeUsecase(contextFactory appcontext.Factory) UpdateUIThemeUsecase {
	return &updateUIThemeUsecase{contextFactory: contextFactory}
}

func (u *updateUIThemeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, bearerToken string, in UpdateUIThemeInput) (*UpdateUIThemeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}

	uiTheme := strings.TrimSpace(in.UITheme)
	if uiTheme == "" {
		uiTheme = "primary"
	}
	if uiTheme != "primary" && uiTheme != "secondary" {
		return nil, apperrors.NewBadRequestError("ui_theme must be primary or secondary")
	}

	if err := app.Repositories.UserProfile.UpdateUITheme(ctx, id, uiTheme); err != nil {
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

	return &UpdateUIThemeOutput{Data: toProfileData(*updated, identity.FullName(info, id), info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

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
		Execute(context.Context, string, bool, UpdateUIThemeInput) (*UpdateUIThemeOutput, apperrors.ApplicationError)
	}

	updateUIThemeUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateUIThemeInput struct {
		ID          string
		UITheme     string `json:"ui_theme"`
		BearerToken string
	}

	UpdateUIThemeOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateUIThemeUsecase(contextFactory appcontext.Factory) UpdateUIThemeUsecase {
	return &updateUIThemeUsecase{contextFactory: contextFactory}
}

func (u *updateUIThemeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input UpdateUIThemeInput) (*UpdateUIThemeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Reuses the assignment-visibility rule: the requester may touch this
	// profile if it's their own, or they administer a school the target
	// belongs to.
	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, input.ID); appErr != nil {
		return nil, appErr
	}

	uiTheme := strings.TrimSpace(input.UITheme)
	if uiTheme == "" {
		uiTheme = "primary"
	}
	if uiTheme != "primary" && uiTheme != "secondary" {
		return nil, apperrors.NewBadRequestError("ui_theme must be primary or secondary")
	}

	if err := app.Repositories.UserProfile.UpdateUITheme(ctx, input.ID, uiTheme); err != nil {
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

	return &UpdateUIThemeOutput{Data: toProfileData(*updated, identity.FullName(info, input.ID), info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

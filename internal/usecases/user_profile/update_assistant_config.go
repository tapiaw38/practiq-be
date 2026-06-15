package userprofile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateAssistantConfigUsecase interface {
		Execute(context.Context, UpdateAssistantConfigInput) (*UpdateAssistantConfigOutput, apperrors.ApplicationError)
	}

	updateAssistantConfigUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateAssistantConfigInput struct {
		ID               string
		AssistantBaseURL string `json:"assistant_base_url"`
		AssistantAPIKey  string `json:"assistant_api_key"`
	}

	UpdateAssistantConfigOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewUpdateAssistantConfigUsecase(contextFactory appcontext.Factory) UpdateAssistantConfigUsecase {
	return &updateAssistantConfigUsecase{contextFactory: contextFactory}
}

func (u *updateAssistantConfigUsecase) Execute(ctx context.Context, input UpdateAssistantConfigInput) (*UpdateAssistantConfigOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	baseURL := strings.TrimSpace(input.AssistantBaseURL)
	apiKey := strings.TrimSpace(input.AssistantAPIKey)

	if err := app.Repositories.UserProfile.UpdateAssistantConfig(ctx, input.ID, baseURL, apiKey); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	return &UpdateAssistantConfigOutput{Data: toProfileData(*updated)}, nil
}

package userprofile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateAssistantConfigUsecase interface {
	Execute(context.Context, UpdateAssistantConfigInput) (*ProfileOutput, apperrors.ApplicationError)
}

type updateAssistantConfigUsecase struct {
	factory appcontext.Factory
}

type UpdateAssistantConfigInput struct {
	ID               string
	AssistantBaseURL string `json:"assistant_base_url"`
	AssistantAPIKey  string `json:"assistant_api_key"`
}

func NewUpdateAssistantConfigUsecase(factory appcontext.Factory) UpdateAssistantConfigUsecase {
	return &updateAssistantConfigUsecase{factory: factory}
}

func (u *updateAssistantConfigUsecase) Execute(ctx context.Context, input UpdateAssistantConfigInput) (*ProfileOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	return &ProfileOutput{Data: toProfileData(*updated)}, nil
}

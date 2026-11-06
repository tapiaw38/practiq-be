package gilliesettings

import (
	"context"
	"strings"

	repo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/gillie_settings"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/secretbox"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, updatedBy string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		BaseURL string `json:"base_url"`
		APIKey  string `json:"api_key"`
	}

	UpdateOutput struct {
		Data SettingsData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, updatedBy string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	baseURL := strings.TrimRight(strings.TrimSpace(in.BaseURL), "/")
	if baseURL == "" {
		return nil, apperrors.NewApplicationError(mappings.GillieSettingsBadURLError, nil)
	}
	if err := assistant.ValidateBaseURL(baseURL); err != nil {
		details := mappings.GillieSettingsBadURLError
		details.Message = err.Error()
		return nil, apperrors.NewApplicationError(details, err)
	}

	settings := repo.Settings{BaseURL: baseURL, UpdatedBy: updatedBy}

	if apiKey := strings.TrimSpace(in.APIKey); apiKey != "" {
		sealed, err := assistantcfg.Seal(apiKey)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.GillieSettingsEncryptionError, err)
		}
		settings.APIKeyEncrypted = sealed
		settings.APIKeyLast4 = secretbox.Last4(apiKey)
	}

	if err := app.Repositories.GillieSettings.Save(ctx, settings); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GillieSettingsSaveError, err)
	}

	stored, err := app.Repositories.GillieSettings.Get(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GillieSettingsGetError, err)
	}
	return &UpdateOutput{Data: toSettingsData(stored)}, nil
}

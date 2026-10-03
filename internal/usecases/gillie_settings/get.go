package gilliesettings

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data SettingsData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	stored, err := app.Repositories.GillieSettings.Get(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GillieSettingsGetError, err)
	}
	return &GetOutput{Data: toSettingsData(stored)}, nil
}

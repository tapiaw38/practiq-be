package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/assistant"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
)

type ProxyUsecase interface {
	Execute(context.Context, ProxyInput) (*assistant.ProxyResponse, apperrors.ApplicationError)
}

type proxyUsecase struct {
	factory appcontext.Factory
}

type ProxyInput struct {
	UserID      string
	Method      string
	Path        string
	ContentType string
	Body        []byte
}

func NewProxyUsecase(factory appcontext.Factory) ProxyUsecase {
	return &proxyUsecase{factory: factory}
}

func (u *proxyUsecase) Execute(ctx context.Context, input ProxyInput) (*assistant.ProxyResponse, apperrors.ApplicationError) {
	app := u.factory()

	profile, err := app.Repositories.UserProfile.Get(ctx, input.UserID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if profile == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	cfg := assistant.Config{
		BaseURL: profile.AssistantBaseURL,
		APIKey:  profile.AssistantAPIKey,
	}
	if !app.AssistantService.IsConfigured(cfg) {
		return nil, apperrors.NewBadRequestError("assistant is not configured for this profile")
	}

	response, proxyErr := app.AssistantService.Proxy(ctx, cfg, input.Method, input.Path, input.ContentType, input.Body)
	if proxyErr != nil {
		return nil, apperrors.NewInternalError(proxyErr)
	}

	return response, nil
}

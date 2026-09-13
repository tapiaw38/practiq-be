package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
)

type (
	ProxyUsecase interface {
		Execute(context.Context, ProxyInput) (*assistant.ProxyResponse, apperrors.ApplicationError)
	}

	proxyUsecase struct {
		contextFactory appcontext.Factory
	}

	ProxyInput struct {
		UserID      string
		Method      string
		Path        string
		ContentType string
		Body        []byte
	}
)

func NewProxyUsecase(contextFactory appcontext.Factory) ProxyUsecase {
	return &proxyUsecase{contextFactory: contextFactory}
}

func (u *proxyUsecase) Execute(ctx context.Context, input ProxyInput) (*assistant.ProxyResponse, apperrors.ApplicationError) {
	app := u.contextFactory()

	cfg := assistantcfg.Resolve(ctx, app)
	if !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return nil, apperrors.NewBadRequestError("the assistant is not configured for this platform")
	}

	response, proxyErr := app.Integrations.AssistantGateway.Proxy(ctx, cfg, input.Method, input.Path, input.ContentType, input.Body)
	if proxyErr != nil {
		return nil, apperrors.NewInternalError(proxyErr)
	}

	return response, nil
}

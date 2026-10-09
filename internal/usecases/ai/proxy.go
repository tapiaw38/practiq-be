package ai

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
)

type (
	ProxyUsecase interface {
		Execute(context.Context, ProxyInput) (*ProxyOutput, apperrors.ApplicationError)
		ExecuteStream(context.Context, ProxyInput, func(int, string, []byte) error) (int, string, apperrors.ApplicationError)
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

	ProxyOutput struct {
		StatusCode  int
		ContentType string
		Body        []byte
	}
)

func NewProxyUsecase(contextFactory appcontext.Factory) ProxyUsecase {
	return &proxyUsecase{contextFactory: contextFactory}
}

func (u *proxyUsecase) Execute(ctx context.Context, input ProxyInput) (*ProxyOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	cfg := assistantcfg.Resolve(ctx, app)
	if !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return nil, apperrors.NewBadRequestError("the assistant is not configured for this platform")
	}

	response, proxyErr := app.Integrations.AssistantGateway.Proxy(ctx, cfg, input.Method, input.Path, input.ContentType, input.Body)
	if proxyErr != nil {
		return nil, apperrors.NewInternalError(proxyErr)
	}

	return toProxyOutput(response), nil
}

func (u *proxyUsecase) ExecuteStream(ctx context.Context, input ProxyInput, write func(int, string, []byte) error) (int, string, apperrors.ApplicationError) {
	app := u.contextFactory()
	cfg := assistantcfg.Resolve(ctx, app)
	if !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return 0, "", apperrors.NewBadRequestError("the assistant is not configured for this platform")
	}
	status, contentType, err := app.Integrations.AssistantGateway.ProxyStream(ctx, cfg, input.Method, input.Path, input.ContentType, input.Body, write)
	if err != nil {
		return 0, "", apperrors.NewInternalError(fmt.Errorf("assistant stream proxy: %w", err))
	}
	return status, contentType, nil
}

func toProxyOutput(response *assistant.ProxyResponse) *ProxyOutput {
	return &ProxyOutput{StatusCode: response.StatusCode, ContentType: response.ContentType, Body: response.Body}
}

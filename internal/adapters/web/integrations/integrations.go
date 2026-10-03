package integrations

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
)

type Integrations struct {
	AssistantGateway assistant.Gateway
	AuthAPI          authapi.Client
	Payments         payments.Client
}

func CreateIntegrations(authAPIURL, paymentsURL, paymentsAPIKey string) *Integrations {
	return &Integrations{
		AssistantGateway: assistant.NewGateway(),
		AuthAPI:          authapi.NewClient(authAPIURL),
		Payments:         payments.NewClient(paymentsURL, paymentsAPIKey),
	}
}

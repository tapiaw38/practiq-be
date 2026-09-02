package integrations

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
)

type Integrations struct {
	AssistantGateway assistant.Gateway
	AuthAPI          authapi.Client
}

func CreateIntegrations(authAPIURL string) *Integrations {
	return &Integrations{
		AssistantGateway: assistant.NewGateway(),
		AuthAPI:          authapi.NewClient(authAPIURL),
	}
}

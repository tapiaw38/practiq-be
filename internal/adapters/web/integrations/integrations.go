package integrations

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
)

type Integrations struct {
	AssistantGateway assistant.Gateway
}

func CreateIntegrations() *Integrations {
	return &Integrations{
		AssistantGateway: assistant.NewGateway(),
	}
}

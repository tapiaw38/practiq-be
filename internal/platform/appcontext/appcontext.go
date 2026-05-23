package appcontext

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-be/internal/platform/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/strategy"
)

type Context struct {
	Repositories     *repositories.Repositories
	KumonStrategy    strategy.LearningStrategyService
	AssistantService assistant.Service
}

type Factory func() *Context

func NewFactory(repos *repositories.Repositories, kumon strategy.LearningStrategyService, assistantService assistant.Service) Factory {
	ctx := &Context{
		Repositories:     repos,
		KumonStrategy:    kumon,
		AssistantService: assistantService,
	}
	return func() *Context {
		return ctx
	}
}

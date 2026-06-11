package appcontext

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-be/internal/platform/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
	"github.com/tapiaw38/practiq-be/internal/platform/strategy"
)

type Context struct {
	Repositories     *repositories.Repositories
	KumonStrategy    strategy.LearningStrategyService
	AssistantService assistant.Service
	ImageStorage     storage.ImageStorage
}

type Factory func() *Context

func NewFactory(repos *repositories.Repositories, kumon strategy.LearningStrategyService, assistantService assistant.Service, imageStorage storage.ImageStorage) Factory {
	ctx := &Context{
		Repositories:     repos,
		KumonStrategy:    kumon,
		AssistantService: assistantService,
		ImageStorage:     imageStorage,
	}
	return func() *Context {
		return ctx
	}
}

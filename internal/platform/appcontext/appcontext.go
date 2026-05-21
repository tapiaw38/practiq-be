package appcontext

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-be/internal/platform/strategy"
)

type Context struct {
	Repositories     *repositories.Repositories
	KumonStrategy    strategy.LearningStrategyService
}

type Factory func() *Context

func NewFactory(repos *repositories.Repositories, kumon strategy.LearningStrategyService) Factory {
	ctx := &Context{
		Repositories:  repos,
		KumonStrategy: kumon,
	}
	return func() *Context {
		return ctx
	}
}

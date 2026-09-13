package appcontext

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

type Context struct {
	Repositories *repositories.Repositories
	Integrations *integrations.Integrations
	ImageStorage storage.ImageStorage
}

type Factory func() *Context

func NewFactory(repos *repositories.Repositories, integ *integrations.Integrations, imageStorage storage.ImageStorage) Factory {
	ctx := &Context{
		Repositories: repos,
		Integrations: integ,
		ImageStorage: imageStorage,
	}
	return func() *Context {
		return ctx
	}
}

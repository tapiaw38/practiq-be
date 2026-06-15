package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, id string) error
	}

	deleteUsecase struct{ contextFactory appcontext.Factory }
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) error {
	app := u.contextFactory()
	return app.Repositories.Notebook.Delete(ctx, id)
}

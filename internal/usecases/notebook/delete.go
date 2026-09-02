package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError
	}

	deleteUsecase struct{ contextFactory appcontext.Factory }
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	// Verify notebook exists and check ownership
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if nb == nil {
		return apperrors.NewNotFoundError("notebook not found")
	}

	if !isSuperAdmin && nb.TeacherID != requesterID {
		return apperrors.NewForbiddenError()
	}

	if err := app.Repositories.Notebook.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.NotebookDeleteError, err)
	}

	return nil
}

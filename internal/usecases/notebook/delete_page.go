package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type DeletePageUsecase interface {
	Execute(ctx context.Context, requesterID string, isSuperAdmin bool, pageID string) apperrors.ApplicationError
}

type deletePageUsecase struct{ contextFactory appcontext.Factory }

func NewDeletePageUsecase(contextFactory appcontext.Factory) DeletePageUsecase {
	return &deletePageUsecase{contextFactory: contextFactory}
}

func (u *deletePageUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, pageID string) apperrors.ApplicationError {
	app := u.contextFactory()
	page, err := app.Repositories.Notebook.GetPage(ctx, pageID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if page == nil {
		return apperrors.NewNotFoundError("notebook page not found")
	}
	notebook, err := app.Repositories.Notebook.Get(ctx, page.NotebookID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if notebook == nil {
		return apperrors.NewNotFoundError("notebook not found")
	}
	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, notebook.CourseID); appErr != nil {
		return appErr
	}
	if err := app.Repositories.Notebook.DeletePage(ctx, pageID); err != nil {
		return apperrors.NewApplicationError(mappings.NotebookDeleteError, err)
	}
	return nil
}

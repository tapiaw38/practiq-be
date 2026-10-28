package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return appErr
	}

	if err := app.Repositories.Course.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.CourseDeleteError, err)
	}

	return nil
}

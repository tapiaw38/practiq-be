package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func requesterCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr == nil {
		return nil
	}
	hasAccess, err := studentHasNotebookCourseAccess(ctx, app, requesterID, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if !hasAccess {
		return apperrors.NewForbiddenError()
	}
	return nil
}

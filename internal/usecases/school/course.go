package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func EnsureCanManageCourse(
	ctx context.Context,
	app *appcontext.Context,
	requesterID string,
	isSuperAdmin bool,
	courseID string,
) (*domain.Course, apperrors.ApplicationError) {
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}
	if isSuperAdmin || course.TeacherID == requesterID {
		return course, nil
	}
	if appErr := EnsureAdministers(ctx, app, requesterID, false, course.SchoolID); appErr != nil {
		return nil, appErr
	}
	return course, nil
}

func EnsureCourseAcceptsWork(ctx context.Context, app *appcontext.Context, courseID string) apperrors.ApplicationError {
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if !course.AcceptsWork() {
		return apperrors.NewBadRequestError("this course is closed: it no longer accepts submissions")
	}
	return nil
}

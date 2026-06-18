package course

import (
	"context"

	reposCourse "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func requesterCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, courseID string) apperrors.ApplicationError {
	if isAdmin {
		return nil
	}

	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if course.TeacherID == requesterID {
		return nil
	}

	courses, err := app.Repositories.Course.List(ctx, reposCourse.ListFilterOptions{StudentID: requesterID})
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseListError, err)
	}
	for _, c := range courses {
		if c.ID == courseID {
			return nil
		}
	}

	return apperrors.NewForbiddenError()
}

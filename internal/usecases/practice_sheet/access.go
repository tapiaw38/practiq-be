package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func requesterCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
	if isSuperAdmin {
		return nil
	}
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseNotFoundError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if course.TeacherID == requesterID {
		return nil
	}
	hasAccess, err := studentHasCourseAccess(ctx, app, requesterID, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseNotFoundError, err)
	}
	if !hasAccess {
		return apperrors.NewForbiddenError()
	}
	return nil
}

// requesterCanViewTeacherData distinguishes a course teacher from a student
// who is allowed to solve the same sheet. It controls answers and canonical
// storage URLs in output, not authorization to read the sheet itself.
func requesterCanViewTeacherData(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) (bool, apperrors.ApplicationError) {
	if isSuperAdmin {
		return true, nil
	}
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return false, apperrors.NewApplicationError(mappings.CourseNotFoundError, err)
	}
	if course == nil {
		return false, apperrors.NewNotFoundError("course not found")
	}
	return course.TeacherID == requesterID, nil
}

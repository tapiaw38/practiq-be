package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// EnsureCanManageCourse resolves a course and refuses anyone who may not change
// it or anything hanging off it — its topics, its exercises, its materials.
//
// Three people may: the teacher who owns it, an admin of the school it belongs
// to, and a platform superadmin. It is the rule Campus already applies through
// campusaccess.CanManageCourse, written here once instead of as a teacher_id
// comparison repeated in every write path — which is why a school's own
// administrator could not touch the courses of the teachers they administer.
//
// Course.Get joins an active school, so a course in a closed school reads as
// not found and no caller, superadmin included, writes to it.
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

// EnsureCourseAcceptsWork refuses a student submission when the course is no
// longer taking any.
//
// Archiving a course is how a teacher says it is finished. Students who took
// it keep reading it — their answers and marks are their own record, and
// deleting that view the day the teacher tidies up would take their history
// with it — but nothing new is accepted, so an archived course cannot quietly
// keep collecting work nobody will grade. A draft refuses for the opposite
// reason: it was never open.
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

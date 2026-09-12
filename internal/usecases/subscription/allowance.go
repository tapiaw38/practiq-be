package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// EnsureCanAddStudent refuses a link that would take a teacher past their plan.
//
// Called wherever a teacher gains a student, so the number on the subscription
// tab is a limit and not a decoration.
//
// A student the teacher already has is always allowed through: the paths that
// call this are idempotent, and re-running one must not start failing because
// the plan happens to be full.
func EnsureCanAddStudent(
	ctx context.Context,
	app *appcontext.Context,
	teacherID, studentID string,
) apperrors.ApplicationError {
	already, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, teacherID, studentID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	if already {
		return nil
	}

	// The limit belongs to a school, not to a teacher. A teacher with their own
	// school who also teaches at an institution must not have those students
	// charged to their personal plan. scopeFor decides all of that, and the
	// subscription screen reads the same answer.
	scope, appErr := scopeFor(ctx, app, teacherID)
	if appErr != nil {
		return appErr
	}
	if !scope.Enforced() {
		return nil
	}

	used, appErr := studentsUsed(ctx, app, scope.SchoolID)
	if appErr != nil {
		return appErr
	}

	if !(domain.TeacherSubscription{Plan: scope.Plan, StudentsUsed: used}).CanAddStudent() {
		return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
	}
	return nil
}

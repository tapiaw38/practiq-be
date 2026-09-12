package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// EnsureCanAddStudent refuses a link that would take a school past its plan.
//
// Called wherever a teacher gains a student, so the number on the subscription
// tab is a limit and not a decoration.
//
// schoolID is the school the student is joining — an invitation's, a course's —
// and it is what the cap is read from. Pass it empty only when there is none to
// name, which means the teacher's own school. Resolving it from the teacher
// instead charged an institution's students against that teacher's personal
// plan and refused them once it filled.
//
// A student who already belongs to the school is always allowed through: the
// paths that call this are idempotent, and re-running one must not start
// failing because the plan happens to be full.
func EnsureCanAddStudent(
	ctx context.Context,
	app *appcontext.Context,
	schoolID, teacherID, studentID string,
) apperrors.ApplicationError {
	// The limit belongs to a school, not to a teacher. A teacher with their own
	// school who also teaches at an institution must not have those students
	// charged to their personal plan. scopeFor decides all of that, and the
	// subscription screen reads the same answer.
	scope, appErr := scopeFor(ctx, app, schoolID, teacherID)
	if appErr != nil {
		return appErr
	}
	if !scope.Enforced() {
		return nil
	}

	// Asked of the same table the count comes from. Asking whether the teacher
	// already had the student instead let one through whose membership was
	// inactive — a downgrade had deactivated them — and the count they were
	// waved past is the one their reactivation then pushed over the plan.
	already, appErr := isActiveMember(ctx, app, scope.SchoolID, studentID)
	if appErr != nil {
		return appErr
	}
	if already {
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

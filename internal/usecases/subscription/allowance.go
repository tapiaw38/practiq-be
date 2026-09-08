package subscription

import (
	"context"
	"log"

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

	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, teacherID)
	if err != nil {
		// Deliberately allowed. Reading the plan failed, so the plan is
		// unknown, and treating unknown as "free plan, one student" would stop
		// paying teachers from working every time the payments service
		// hiccups. Letting one extra student in costs a little revenue;
		// refusing a teacher their class costs the product.
		log.Printf("[payments] allowance check skipped teacher_id=%s err=%v", teacherID, err)
		return nil
	}

	plan := domain.FreePlan
	if entitlement != nil && entitlement.Active {
		plan = domain.PlanFromMetadata(0, "", entitlement.Metadata)
	}

	used, err := app.Repositories.TeacherStudentAssignment.CountStudents(ctx, teacherID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}

	if !(domain.TeacherSubscription{Plan: plan, StudentsUsed: used}).CanAddStudent() {
		return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
	}
	return nil
}

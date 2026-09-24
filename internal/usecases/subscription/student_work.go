package subscription

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// StudentsCanWork reports whether a school's students may still add work.
//
// Reading is never in question: whatever they already wrote stays theirs. This
// is about submitting practices and notebook pages, which is what a school
// stops being paid for.
//
// Answers true whenever the answer is not clearly no. A payments outage, an
// institution invoiced outside the product, a trial still running — none of
// those are somebody failing to pay, and refusing a class over any of them
// costs far more than the revenue it protects.
func StudentsCanWork(ctx context.Context, app *appcontext.Context, schoolID, ownerID string) bool {
	scope, appErr := scopeFor(ctx, app, schoolID, ownerID)
	if appErr != nil {
		return true
	}
	switch scope.State {
	case capNone, capUnknown:
		return true
	}
	if scope.Active || scope.GraceEndsAt != nil {
		return true
	}
	// Left with the free month. It runs out to an allowance of zero, and a
	// school nobody is paying for and whose trial is spent is the one case
	// where the work stops.
	return scope.Plan.MaxStudents > 0
}

// StudentMayWork narrows StudentsCanWork to one student.
//
// A school can be paid up and still hold more students than its plan allows:
// changing to a smaller plan does not deactivate anybody, deliberately, so the
// teacher gets to choose who stays. Until they do, the cap is applied in the
// order the choice would default to — otherwise a teacher on a plan of five
// keeps fifteen students working by never opening the screen.
//
// The teacher's decision, once made, is stored on the membership and wins:
// this only ever runs for somebody still marked active.
func StudentMayWork(ctx context.Context, app *appcontext.Context, schoolID, ownerID, studentID string) bool {
	if !StudentsCanWork(ctx, app, schoolID, ownerID) {
		return false
	}
	scope, appErr := scopeFor(ctx, app, schoolID, ownerID)
	if appErr != nil || !scope.Enforced() || scope.SchoolID == "" {
		return true
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, scope.SchoolID)
	if err != nil {
		log.Printf("[subscription] activity lookup failed school_id=%s err=%v", scope.SchoolID, err)
		return true
	}
	for _, id := range domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, nil) {
		if id == studentID {
			return false
		}
	}
	return true
}

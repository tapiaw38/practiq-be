package subscription

import (
	"context"

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

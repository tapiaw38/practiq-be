package subscription

import (
	"context"
	"log"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// capState says why a limit does or does not apply.
//
// The three cases look the same to the check that adds a student — none of
// them refuse anyone — but they do not look the same on the subscription
// screen, and collapsing them is how it ends up telling a paying teacher they
// have no limit during an outage.
type capState int

const (
	// capEnforced is a real allowance, read from the plan the teacher pays for
	// or from the free month.
	capEnforced capState = iota
	// capNone is a teacher no per-student limit applies to: one who teaches at
	// an institution, or whose school is invoiced outside the product.
	capNone
	// capUnknown is a payments outage. Nothing is enforced, and the screen
	// falls back to the free plan rather than claiming there is no limit.
	capUnknown
)

// planScope is the subscription rules resolved for one teacher: which school
// their students are counted in, and what their plan allows.
type planScope struct {
	// SchoolID is where students are counted. Empty when the teacher has no
	// school of their own, whose students are on nobody's subscription.
	SchoolID string
	Plan     domain.TeacherPlan
	State    capState
	// Active is a teacher paying right now — not merely one whose free month
	// has yet to run out.
	Active bool
	// RenewsAt is when the paid period ends. Nil on the free plan, where the
	// date that matters is the end of the trial instead.
	RenewsAt *time.Time
}

// Enforced reports whether a limit should be applied at all.
func (s planScope) Enforced() bool { return s.State == capEnforced }

// scopeFor resolves what a teacher's subscription allows right now.
//
// Every rule about who is capped and by how much lives here. It used to be
// written once per consumer, and the copies drifted: the screen counted a
// teacher's students by assignment while the limit counted a school's members,
// so the two disagreed about whether the plan was full.
func scopeFor(ctx context.Context, app *appcontext.Context, teacherID string) (planScope, apperrors.ApplicationError) {
	school, err := app.Repositories.School.GetPersonal(ctx, teacherID)
	if err != nil {
		return planScope{}, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	// Only teachers at institutions have no school of their own, and an
	// institution's students are not on anybody's subscription.
	if school == nil || school.Billing == domain.SchoolBillingDirect {
		return planScope{State: capNone}, nil
	}

	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, teacherID)
	if err != nil {
		// Deliberately allowed. Reading the plan failed, so the plan is
		// unknown, and treating unknown as "free plan" would stop paying
		// teachers from working every time the payments service hiccups.
		// Letting one extra student in costs a little revenue; refusing a
		// teacher their class costs the product.
		log.Printf("[payments] plan lookup failed teacher_id=%s err=%v", teacherID, err)
		return planScope{SchoolID: school.ID, Plan: domain.FreePlan, State: capUnknown}, nil
	}

	if entitlement != nil && entitlement.Active {
		planID := 0
		if entitlement.PlanID != nil {
			planID = *entitlement.PlanID
		}
		return planScope{
			SchoolID: school.ID,
			Plan:     domain.PlanFromMetadata(planID, planName(entitlement.Metadata), entitlement.Metadata),
			State:    capEnforced,
			Active:   true,
			RenewsAt: entitlement.AccessUntil,
		}, nil
	}

	// Nobody paid, so the free month decides — and it is the profile's age
	// that says whether there is any of it left.
	profile, err := app.Repositories.UserProfile.Get(ctx, teacherID)
	if err != nil {
		return planScope{}, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return planScope{}, apperrors.NewNotFoundError("profile not found")
	}

	return planScope{
		SchoolID: school.ID,
		Plan:     domain.EffectiveFreePlan(profile.CreatedAt, time.Now().UTC()),
		State:    capEnforced,
	}, nil
}

// studentsUsed counts a school's active students: the same number the limit is
// checked against, wherever it is shown.
func studentsUsed(ctx context.Context, app *appcontext.Context, schoolID string) (int, apperrors.ApplicationError) {
	if schoolID == "" {
		return 0, nil
	}
	used, err := app.Repositories.School.CountStudents(ctx, schoolID)
	if err != nil {
		return 0, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return used, nil
}

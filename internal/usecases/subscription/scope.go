package subscription

import (
	"context"
	"log"
	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
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
	// Status is what the gateway calls the agreement. An Active teacher may
	// still be paused — they keep the month they bought — and that has to
	// reach the screen or there is nothing to resume.
	Status string
	// GraceEndsAt is set only while a lapsed plan is still being honoured. The
	// teacher keeps everything until then, and has to be told the date.
	GraceEndsAt *time.Time
}

// Enforced reports whether a limit should be applied at all.
func (s planScope) Enforced() bool { return s.State == capEnforced }

// scopeFor resolves what a plan allows for students joining a given school.
//
// The school is the one the student ends up in, not the one the teacher
// happens to own. Every teacher gets a personal school on sign-up
// (ensurePersonalSchool), so resolving the cap from the teacher charged an
// institution's students against that teacher's own plan — and refused them
// once it filled, even though institutions are invoiced by contract and have
// no student limit at all.
//
// An empty schoolID means the teacher's own school: the only caller that has
// no school to name is a direct teacher-student assignment, which is exactly
// the personal case.
//
// Every rule about who is capped and by how much lives here. It used to be
// written once per consumer, and the copies drifted: the screen counted a
// teacher's students by assignment while the limit counted a school's members,
// so the two disagreed about whether the plan was full.
func scopeFor(ctx context.Context, app *appcontext.Context, schoolID, teacherID string) (planScope, apperrors.ApplicationError) {
	school, appErr := resolveSchool(ctx, app, schoolID, teacherID)
	if appErr != nil {
		return planScope{}, appErr
	}
	if school == nil {
		return planScope{State: capNone}, nil
	}

	// An institution is invoiced by contract, so nothing about it is capped
	// per student — whatever its billing column says. Its own admins decide how
	// many students, teachers and admins it has, and none of that is sold here.
	if school.Kind == domain.SchoolKindInstitution || school.Billing == domain.SchoolBillingDirect {
		return planScope{SchoolID: school.ID, State: capNone}, nil
	}

	// A personal school is its owner's, and it is their plan that pays for it —
	// not the plan of whoever happens to be adding the student.
	owner := school.CreatedBy
	if owner == "" {
		owner = teacherID
	}

	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, owner)
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
			RenewsAt: renewsAt(entitlement.AccessUntil),
			Status:   entitlement.Status,
		}, nil
	}

	// Paid until recently still counts. A card that expires on a Friday should
	// not lock a class out on Saturday, so the plan is honoured in full for a
	// few days more and the teacher is told when that ends.
	if grace, found := graceScope(ctx, app, owner, school.ID); found {
		return grace, nil
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

// resolveSchool reads the named school, or the teacher's own when none is
// named. A nil school is one that does not exist, which caps nothing: there is
// no plan behind a school that is not there.
func resolveSchool(ctx context.Context, app *appcontext.Context, schoolID, teacherID string) (*domain.School, apperrors.ApplicationError) {
	var (
		school *domain.School
		err    error
	)
	if schoolID != "" {
		school, err = app.Repositories.School.Get(ctx, schoolID)
	} else {
		school, err = app.Repositories.School.GetPersonal(ctx, teacherID)
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return school, nil
}

// isActiveMember reports whether a student already counts towards this
// school's total, which is what makes adding them again free.
func isActiveMember(ctx context.Context, app *appcontext.Context, schoolID, studentID string) (bool, apperrors.ApplicationError) {
	if schoolID == "" {
		return false, nil
	}
	members, err := app.Repositories.School.ListForUser(ctx, studentID)
	if err != nil {
		return false, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.SchoolID == schoolID && member.Active {
			return true, nil
		}
	}
	return false, nil
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

// renewsAt unwraps the payments timestamp, which tolerates a missing timezone.
func renewsAt(access *payments.Timestamp) *time.Time {
	if access == nil || access.IsZero() {
		return nil
	}
	at := access.Time
	return &at
}

// graceScope honours a plan whose paid period ended within domain.GraceDays.
//
// Read from the subscriptions rather than the entitlement on purpose: the
// entitlement answers "is this paid for", and the honest answer here is no.
// What is true is that it was paid for until very recently.
func graceScope(ctx context.Context, app *appcontext.Context, teacherID, schoolID string) (planScope, bool) {
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {
		// Same rule as every other payments failure here: unknown is not
		// "free plan", and an outage must not shrink somebody's school.
		log.Printf("[payments] grace lookup failed teacher_id=%s err=%v", teacherID, err)
		return planScope{}, false
	}

	now := time.Now().UTC()
	var lapsed *payments.Subscription
	for i := range subscriptions {
		candidate := subscriptions[i]
		if candidate.CurrentPeriodEnd == nil {
			continue
		}
		if !domain.InGrace(candidate.CurrentPeriodEnd.UTC(), now) {
			continue
		}
		// The one that was paid furthest into the future is the one they were
		// actually on; older agreements from plan changes are not it.
		if lapsed == nil || candidate.CurrentPeriodEnd.After(lapsed.CurrentPeriodEnd.Time) {
			lapsed = &candidate
		}
	}
	if lapsed == nil {
		return planScope{}, false
	}

	plans, err := app.Integrations.Payments.ListPlans(ctx)
	if err != nil {
		log.Printf("[payments] grace plan lookup failed teacher_id=%s err=%v", teacherID, err)
		return planScope{}, false
	}
	for _, plan := range plans {
		if plan.ID != lapsed.PlanID {
			continue
		}
		endsAt := domain.GraceEndsAt(lapsed.CurrentPeriodEnd.UTC())
		return planScope{
			SchoolID: schoolID,
			Plan:     domain.PlanFromMetadata(plan.ID, planName(plan.Metadata), plan.Metadata),
			State:    capEnforced,
			// Not Active: nothing is being charged, and the screen must not
			// show this as a running subscription.
			Active:      false,
			Status:      lapsed.Status,
			GraceEndsAt: &endsAt,
		}, true
	}
	return planScope{}, false
}

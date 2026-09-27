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

type capState int

const (
	capEnforced capState = iota

	capNone

	capUnknown
)

type planScope struct {
	SchoolID string
	Plan     domain.TeacherPlan
	State    capState

	Active bool

	RenewsAt *time.Time

	Status string

	GraceEndsAt *time.Time
}

func (s planScope) Enforced() bool { return s.State == capEnforced }

func scopeFor(ctx context.Context, app *appcontext.Context, schoolID, teacherID string) (planScope, apperrors.ApplicationError) {
	school, appErr := resolveSchool(ctx, app, schoolID, teacherID)
	if appErr != nil {
		return planScope{}, appErr
	}
	if school == nil {
		return planScope{State: capNone}, nil
	}

	if school.Kind == domain.SchoolKindInstitution || school.Billing == domain.SchoolBillingDirect {
		return planScope{SchoolID: school.ID, State: capNone}, nil
	}

	owner := school.CreatedBy
	if owner == "" {
		owner = teacherID
	}

	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, owner)
	if err != nil {

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

	if grace, found := graceScope(ctx, app, owner, school.ID); found {
		return grace, nil
	}

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

func renewsAt(access *payments.Timestamp) *time.Time {
	if access == nil || access.IsZero() {
		return nil
	}
	at := access.Time
	return &at
}

func graceScope(ctx context.Context, app *appcontext.Context, teacherID, schoolID string) (planScope, bool) {
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {

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

			Active:      false,
			Status:      lapsed.Status,
			GraceEndsAt: &endsAt,
		}, true
	}
	return planScope{}, false
}

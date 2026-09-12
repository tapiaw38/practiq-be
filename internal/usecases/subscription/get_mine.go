package subscription

import (
	"context"
	"log"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
)

type (
	// GetMineUsecase answers what the asking teacher's plan is and how much of
	// it they are using.
	//
	// It reads the paid side from the payments service and the used side from
	// here, because only this side knows who counts as a teacher's student.
	GetMineUsecase interface {
		Execute(ctx context.Context, teacherID string) (*GetMineOutput, apperrors.ApplicationError)
	}

	getMineUsecase struct {
		contextFactory appcontext.Factory
	}

	PlanData struct {
		// PlanID is absent on the free plan, which the payments service does
		// not know about.
		PlanID      int    `json:"plan_id,omitempty"`
		Name        string `json:"name"`
		MaxStudents int    `json:"max_students"`
	}

	SubscriptionData struct {
		Plan   PlanData `json:"plan"`
		Active bool     `json:"active"`
		// Status is the gateway's word for it: authorized, paused, cancelled,
		// or empty on the free plan. Needed because a paused subscription is
		// not an entitlement, so "not active" alone cannot tell a teacher who
		// paused from one who never paid.
		Status        string `json:"status,omitempty"`
		StudentsUsed  int    `json:"students_used"`
		CanAddStudent bool   `json:"can_add_student"`
		RenewsAt      string `json:"renews_at,omitempty"`
		// Uncapped is a teacher no student limit applies to: one who teaches at
		// an institution, or whose school is invoiced outside the product.
		// Without it the screen would show them a plan they are not on and a
		// limit they do not have.
		Uncapped bool `json:"uncapped,omitempty"`
		// TrialExpired is a free month that has run out. The allowance is
		// already zero, and this says why it is rather than leaving a teacher
		// to work out that "0 de 0" means they have to subscribe.
		TrialExpired bool `json:"trial_expired,omitempty"`
	}

	GetMineOutput struct {
		Data SubscriptionData `json:"data"`
	}
)

func NewGetMineUsecase(contextFactory appcontext.Factory) GetMineUsecase {
	return &getMineUsecase{contextFactory: contextFactory}
}

func (u *getMineUsecase) Execute(ctx context.Context, teacherID string) (*GetMineOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// The same resolver the limit check uses, so the number on this screen and
	// the number that refuses a student are one number. They were two, read
	// from different tables, and disagreed whenever a downgrade deactivated
	// somebody or a school membership failed to be written.
	scope, appErr := scopeFor(ctx, app, teacherID)
	if appErr != nil {
		return nil, appErr
	}
	used, appErr := studentsUsed(ctx, app, scope.SchoolID)
	if appErr != nil {
		return nil, appErr
	}

	data := toData(domain.TeacherSubscription{
		Plan:         scope.Plan,
		Active:       scope.Active,
		StudentsUsed: used,
		RenewsAt:     scope.RenewsAt,
	})

	switch scope.State {
	case capNone:
		// No plan to show and nothing to measure against.
		data.Uncapped = true
		data.CanAddStudent = true
		data.Plan = PlanData{Name: "Sin límite"}
	case capUnknown:
		// A payments outage leaves the teacher on the free plan for the length
		// of the outage rather than breaking the screen, and never refuses a
		// student over a number it could not read.
		data.CanAddStudent = true
	}

	// The free month's end is the date that matters to a teacher who is not
	// paying, and it is the only one they have.
	if !scope.Active && !data.Uncapped {
		data.TrialExpired = scope.Plan.MaxStudents == 0
		data.RenewsAt = trialEndsAt(ctx, app, teacherID)
	}

	data.Status = subscriptionStatus(ctx, app, teacherID, scope.Active)
	return &GetMineOutput{Data: data}, nil
}

// trialEndsAt is when the free month runs out, for a teacher who is on it.
func trialEndsAt(ctx context.Context, app *appcontext.Context, teacherID string) string {
	profile, err := app.Repositories.UserProfile.Get(ctx, teacherID)
	if err != nil || profile == nil {
		log.Printf("[subscription] trial date unavailable teacher_id=%s err=%v", teacherID, err)
		return ""
	}
	return domain.TrialEndsAt(profile.CreatedAt).UTC().Format(time.RFC3339)
}

// subscriptionStatus reports the gateway's status for a teacher who is not
// entitled, which is the only way a paused subscription is visible at all.
//
// Without it a teacher who pauses sees the free plan and no way back: pausing
// removes the entitlement, and the entitlement is all the screen would know.
func subscriptionStatus(ctx context.Context, app *appcontext.Context, teacherID string, active bool) string {
	if active {
		return "authorized"
	}
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {
		log.Printf("[payments] subscription lookup failed teacher_id=%s err=%v", teacherID, err)
		return ""
	}
	for _, subscription := range subscriptions {
		if subscription.Status == "paused" {
			return "paused"
		}
	}
	return ""
}

// planName reads the display name a plan chose to publish, falling back to a
// neutral one: the entitlement carries metadata, not the plan's own name.
func planName(metadata map[string]any) string {
	if name, ok := metadata["name"].(string); ok && name != "" {
		return name
	}
	return "Plan activo"
}

func toData(s domain.TeacherSubscription) SubscriptionData {
	data := SubscriptionData{
		Plan: PlanData{
			PlanID:      s.Plan.PlanID,
			Name:        s.Plan.Name,
			MaxStudents: s.Plan.MaxStudents,
		},
		Active:        s.Active,
		StudentsUsed:  s.StudentsUsed,
		CanAddStudent: s.CanAddStudent(),
	}
	if s.RenewsAt != nil {
		data.RenewsAt = s.RenewsAt.UTC().Format(time.RFC3339)
	}
	return data
}

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
	GetMineUsecase interface {
		Execute(ctx context.Context, teacherID string) (*GetMineOutput, apperrors.ApplicationError)
	}

	getMineUsecase struct {
		contextFactory appcontext.Factory
	}

	PlanData struct {
		PlanID      int    `json:"plan_id,omitempty"`
		Name        string `json:"name"`
		MaxStudents int    `json:"max_students"`
	}

	SubscriptionData struct {
		Plan   PlanData `json:"plan"`
		Active bool     `json:"active"`

		Status        string `json:"status,omitempty"`
		StudentsUsed  int    `json:"students_used"`
		CanAddStudent bool   `json:"can_add_student"`
		RenewsAt      string `json:"renews_at,omitempty"`

		Uncapped bool `json:"uncapped,omitempty"`

		GraceEndsAt string `json:"grace_ends_at,omitempty"`

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

	scope, appErr := scopeFor(ctx, app, "", teacherID)
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

		data.Uncapped = true
		data.CanAddStudent = true
		data.Plan = PlanData{Name: "Sin límite"}
	case capUnknown:

		data.CanAddStudent = true
	}

	if scope.GraceEndsAt != nil {

		data.GraceEndsAt = scope.GraceEndsAt.UTC().Format(time.RFC3339)
		data.RenewsAt = ""
	} else if !scope.Active && !data.Uncapped {

		data.TrialExpired = scope.Plan.MaxStudents == 0
		data.RenewsAt = trialEndsAt(ctx, app, teacherID)
	}

	data.Status = subscriptionStatus(ctx, app, teacherID, scope)
	return &GetMineOutput{Data: data}, nil
}

func trialEndsAt(ctx context.Context, app *appcontext.Context, teacherID string) string {
	profile, err := app.Repositories.UserProfile.Get(ctx, teacherID)
	if err != nil || profile == nil {
		log.Printf("[subscription] trial date unavailable teacher_id=%s err=%v", teacherID, err)
		return ""
	}
	return domain.TrialEndsAt(profile.CreatedAt).UTC().Format(time.RFC3339)
}

func subscriptionStatus(ctx context.Context, app *appcontext.Context, teacherID string, scope planScope) string {
	if scope.Active {

		if scope.Status != "" {
			return scope.Status
		}
		return "authorized"
	}
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {
		log.Printf("[payments] subscription lookup failed teacher_id=%s err=%v", teacherID, err)
		return ""
	}

	pending := false
	for _, subscription := range subscriptions {
		switch subscription.Status {
		case "paused":
			return "paused"
		case "pending":
			pending = true
		}
	}
	if pending {
		return "pending"
	}
	return ""
}

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

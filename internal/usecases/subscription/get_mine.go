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
		Plan          PlanData `json:"plan"`
		Active        bool     `json:"active"`
		StudentsUsed  int      `json:"students_used"`
		CanAddStudent bool     `json:"can_add_student"`
		RenewsAt      string   `json:"renews_at,omitempty"`
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

	used, err := app.Repositories.TeacherStudentAssignment.CountStudents(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}

	profile, err := app.Repositories.UserProfile.Get(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewNotFoundError("profile not found")
	}

	subscription := domain.FreeSubscription(profile.CreatedAt, used)

	// A payments outage leaves the teacher on the free plan for the length of
	// the outage rather than breaking the screen. It is the safe direction:
	// showing less than they paid for is recoverable, and it never hands out
	// an allowance nobody paid for.
	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, teacherID)
	if err != nil {
		log.Printf("[payments] entitlement lookup failed teacher_id=%s err=%v", teacherID, err)
	} else if entitlement != nil && entitlement.Active {
		planID := 0
		if entitlement.PlanID != nil {
			planID = *entitlement.PlanID
		}
		subscription = domain.TeacherSubscription{
			Plan:         domain.PlanFromMetadata(planID, planName(entitlement.Metadata), entitlement.Metadata),
			Active:       true,
			StudentsUsed: used,
			RenewsAt:     entitlement.AccessUntil,
		}
	}

	return &GetMineOutput{Data: toData(subscription)}, nil
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

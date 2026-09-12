package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	// PlansUsecase is the plan catalogue. Reading it is open to any teacher —
	// they need to see what they could move to — while changing it is not.
	PlansUsecase interface {
		List(ctx context.Context) (*PlansOutput, apperrors.ApplicationError)
		Create(ctx context.Context, in PlanInput) (*PlanOutput, apperrors.ApplicationError)
		Update(ctx context.Context, planID int, in PlanInput) (*PlanOutput, apperrors.ApplicationError)
		// Deactivate takes a plan off the shelf. Nothing deletes a plan:
		// subscriptions point at theirs, and people keep paying for what they
		// bought after it stops being sold.
		Deactivate(ctx context.Context, planID int) (*PlanOutput, apperrors.ApplicationError)
	}

	plansUsecase struct {
		contextFactory appcontext.Factory
	}

	// PlanInput is what a superadmin sets. MaxStudents is Practiq's own idea
	// of what the plan grants and travels in the payments plan's metadata,
	// which that service stores without reading.
	PlanInput struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Amount      *float64 `json:"amount"`
		Currency    string   `json:"currency"`
		Interval    string   `json:"interval"`
		MaxStudents *int     `json:"max_students"`
		Active      *bool    `json:"active"`
	}

	CatalogPlanData struct {
		PlanID      int     `json:"plan_id"`
		Name        string  `json:"name"`
		Description string  `json:"description,omitempty"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Interval    string  `json:"interval"`
		MaxStudents int     `json:"max_students"`
		Active      bool    `json:"active"`
	}

	PlansOutput struct {
		Data []CatalogPlanData `json:"data"`
	}

	PlanOutput struct {
		Data CatalogPlanData `json:"data"`
	}
)

func NewPlansUsecase(contextFactory appcontext.Factory) PlansUsecase {
	return &plansUsecase{contextFactory: contextFactory}
}

func (u *plansUsecase) List(ctx context.Context) (*PlansOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plans, err := app.Integrations.Payments.ListPlans(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}

	data := make([]CatalogPlanData, 0, len(plans))
	for _, plan := range plans {
		data = append(data, toCatalogPlan(plan))
	}
	return &PlansOutput{Data: data}, nil
}

func (u *plansUsecase) Create(ctx context.Context, in PlanInput) (*PlanOutput, apperrors.ApplicationError) {
	return u.write(ctx, func(client payments.Client) (*payments.Plan, error) {
		return client.CreatePlan(ctx, toPlanInput(in))
	})
}

func (u *plansUsecase) Update(ctx context.Context, planID int, in PlanInput) (*PlanOutput, apperrors.ApplicationError) {
	return u.write(ctx, func(client payments.Client) (*payments.Plan, error) {
		return client.UpdatePlan(ctx, planID, toPlanInput(in))
	})
}

func (u *plansUsecase) Deactivate(ctx context.Context, planID int) (*PlanOutput, apperrors.ApplicationError) {
	return u.write(ctx, func(client payments.Client) (*payments.Plan, error) {
		return client.DeactivatePlan(ctx, planID)
	})
}

func (u *plansUsecase) write(
	ctx context.Context,
	call func(payments.Client) (*payments.Plan, error),
) (*PlanOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plan, err := call(app.Integrations.Payments)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return &PlanOutput{Data: toCatalogPlan(*plan)}, nil
}

// toPlanInput puts Practiq's allowance into the plan's metadata. Only what the
// caller named is sent, so editing a price cannot blank a limit by omission.
func toPlanInput(in PlanInput) payments.PlanInput {
	out := payments.PlanInput{
		Name:        in.Name,
		Description: in.Description,
		Amount:      in.Amount,
		Currency:    in.Currency,
		Interval:    in.Interval,
		Active:      in.Active,
	}
	if in.MaxStudents != nil {
		out.Metadata = map[string]any{
			"max_students": *in.MaxStudents,
			// Carried so the teacher's own screen can name their plan without
			// reading the catalogue.
			"name": in.Name,
		}
	}
	return out
}

func toCatalogPlan(plan payments.Plan) CatalogPlanData {
	return CatalogPlanData{
		PlanID:      plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
		Amount:      plan.Amount,
		Currency:    plan.Currency,
		Interval:    plan.Interval,
		MaxStudents: domain.PlanFromMetadata(plan.ID, plan.Name, plan.Metadata).MaxStudents,
		Active:      plan.Active == 1,
	}
}

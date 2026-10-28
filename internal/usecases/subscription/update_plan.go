package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdatePlanUsecase interface {
		Execute(ctx context.Context, planID int, in UpdatePlanInput) (*UpdatePlanOutput, apperrors.ApplicationError)
	}

	updatePlanUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdatePlanInput struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Amount      *float64 `json:"amount"`
		Currency    string   `json:"currency"`
		Interval    string   `json:"interval"`
		MaxStudents *int     `json:"max_students"`
		Active      *bool    `json:"active"`
	}

	UpdatePlanOutput struct {
		Data CatalogPlanData `json:"data"`
	}
)

func NewUpdatePlanUsecase(contextFactory appcontext.Factory) UpdatePlanUsecase {
	return &updatePlanUsecase{contextFactory: contextFactory}
}

func (u *updatePlanUsecase) Execute(ctx context.Context, planID int, in UpdatePlanInput) (*UpdatePlanOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plan, err := app.Integrations.Payments.UpdatePlan(ctx, planID, toUpdatePlanInput(in))
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return &UpdatePlanOutput{Data: toCatalogPlan(*plan)}, nil
}

func toUpdatePlanInput(in UpdatePlanInput) payments.PlanInput {
	out := payments.PlanInput{Name: in.Name, Description: in.Description, Amount: in.Amount, Currency: in.Currency, Interval: in.Interval, Active: in.Active}
	if in.MaxStudents != nil {
		out.Metadata = map[string]any{"max_students": *in.MaxStudents, "name": in.Name}
	}
	return out
}

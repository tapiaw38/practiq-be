package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreatePlanUsecase interface {
		Execute(ctx context.Context, in CreatePlanInput) (*CreatePlanOutput, apperrors.ApplicationError)
	}

	createPlanUsecase struct {
		contextFactory appcontext.Factory
	}

	CreatePlanInput struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Amount      *float64 `json:"amount"`
		Currency    string   `json:"currency"`
		Interval    string   `json:"interval"`
		MaxStudents *int     `json:"max_students"`
		Active      *bool    `json:"active"`
	}

	CreatePlanOutput struct {
		Data CatalogPlanData `json:"data"`
	}
)

func NewCreatePlanUsecase(contextFactory appcontext.Factory) CreatePlanUsecase {
	return &createPlanUsecase{contextFactory: contextFactory}
}

func (u *createPlanUsecase) Execute(ctx context.Context, in CreatePlanInput) (*CreatePlanOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plan, err := app.Integrations.Payments.CreatePlan(ctx, toCreatePlanInput(in))
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return &CreatePlanOutput{Data: toCatalogPlan(*plan)}, nil
}

func toCreatePlanInput(in CreatePlanInput) payments.PlanInput {
	out := payments.PlanInput{Name: in.Name, Description: in.Description, Amount: in.Amount, Currency: in.Currency, Interval: in.Interval, Active: in.Active}
	if in.MaxStudents != nil {
		out.Metadata = map[string]any{"max_students": *in.MaxStudents, "name": in.Name}
	}
	return out
}

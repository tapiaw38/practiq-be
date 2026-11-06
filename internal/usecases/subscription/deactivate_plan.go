package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DeactivatePlanUsecase interface {
		Execute(ctx context.Context, planID int) (*DeactivatePlanOutput, apperrors.ApplicationError)
	}

	deactivatePlanUsecase struct {
		contextFactory appcontext.Factory
	}

	DeactivatePlanOutput struct {
		Data CatalogPlanData `json:"data"`
	}
)

func NewDeactivatePlanUsecase(contextFactory appcontext.Factory) DeactivatePlanUsecase {
	return &deactivatePlanUsecase{contextFactory: contextFactory}
}

func (u *deactivatePlanUsecase) Execute(ctx context.Context, planID int) (*DeactivatePlanOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plan, err := app.Integrations.Payments.DeactivatePlan(ctx, planID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return &DeactivatePlanOutput{Data: toCatalogPlan(*plan)}, nil
}

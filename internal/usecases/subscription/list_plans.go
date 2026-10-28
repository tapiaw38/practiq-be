package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListPlansUsecase interface {
		Execute(ctx context.Context) (*ListPlansOutput, apperrors.ApplicationError)
	}

	listPlansUsecase struct {
		contextFactory appcontext.Factory
	}

	ListPlansOutput struct {
		Data []CatalogPlanData `json:"data"`
	}
)

func NewListPlansUsecase(contextFactory appcontext.Factory) ListPlansUsecase {
	return &listPlansUsecase{contextFactory: contextFactory}
}

func (u *listPlansUsecase) Execute(ctx context.Context) (*ListPlansOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	plans, err := app.Integrations.Payments.ListPlans(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}

	data := make([]CatalogPlanData, 0, len(plans))
	for _, plan := range plans {
		data = append(data, toCatalogPlan(plan))
	}
	return &ListPlansOutput{Data: data}, nil
}

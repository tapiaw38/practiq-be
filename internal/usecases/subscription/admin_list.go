package subscription

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	AdminListUsecase interface {
		Execute(ctx context.Context, status string) (*AdminListOutput, apperrors.ApplicationError)
	}

	adminListUsecase struct {
		contextFactory appcontext.Factory
	}

	AdminBillingCycleData struct {
		Status           string  `json:"status"`
		Amount           float64 `json:"amount"`
		Currency         string  `json:"currency"`
		PeriodStart      string  `json:"period_start"`
		PeriodEnd        string  `json:"period_end"`
		GatewayPaymentID string  `json:"gateway_payment_id,omitempty"`
	}

	AdminSubscriptionData struct {
		ID               int                    `json:"id"`
		UserID           string                 `json:"user_id"`
		PayerEmail       string                 `json:"payer_email,omitempty"`
		PlanID           int                    `json:"plan_id"`
		PlanName         string                 `json:"plan_name"`
		Status           string                 `json:"status"`
		CurrentPeriodEnd string                 `json:"current_period_end,omitempty"`
		NextPaymentAt    string                 `json:"next_payment_at,omitempty"`
		CreatedAt        string                 `json:"created_at"`
		LastCycle        *AdminBillingCycleData `json:"last_cycle,omitempty"`
	}

	AdminListOutput struct {
		Data []AdminSubscriptionData `json:"data"`
	}
)

func NewAdminListUsecase(contextFactory appcontext.Factory) AdminListUsecase {
	return &adminListUsecase{contextFactory: contextFactory}
}

func (u *adminListUsecase) Execute(ctx context.Context, status string) (*AdminListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	subscriptions, err := app.Integrations.Payments.ListAllSubscriptions(ctx, status)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}

	data := make([]AdminSubscriptionData, 0, len(subscriptions))
	for _, sub := range subscriptions {
		data = append(data, toAdminSubscriptionData(sub))
	}
	return &AdminListOutput{Data: data}, nil
}

func toAdminSubscriptionData(sub payments.AdminSubscription) AdminSubscriptionData {
	data := AdminSubscriptionData{
		ID:         sub.ID,
		UserID:     sub.UserID,
		PayerEmail: sub.PayerEmail,
		PlanID:     sub.PlanID,
		PlanName:   sub.PlanName,
		Status:     sub.Status,
		CreatedAt:  sub.CreatedAt.UTC().Format(time.RFC3339),
	}
	if sub.CurrentPeriodEnd != nil {
		data.CurrentPeriodEnd = sub.CurrentPeriodEnd.UTC().Format(time.RFC3339)
	}
	if sub.NextPaymentAt != nil {
		data.NextPaymentAt = sub.NextPaymentAt.UTC().Format(time.RFC3339)
	}
	if sub.LastCycle != nil {
		cycle := &AdminBillingCycleData{
			Status:           sub.LastCycle.Status,
			Amount:           sub.LastCycle.Amount,
			Currency:         sub.LastCycle.Currency,
			GatewayPaymentID: sub.LastCycle.GatewayPaymentID,
		}
		if sub.LastCycle.PeriodStart != nil {
			cycle.PeriodStart = sub.LastCycle.PeriodStart.UTC().Format(time.RFC3339)
		}
		if sub.LastCycle.PeriodEnd != nil {
			cycle.PeriodEnd = sub.LastCycle.PeriodEnd.UTC().Format(time.RFC3339)
		}
		data.LastCycle = cycle
	}
	return data
}

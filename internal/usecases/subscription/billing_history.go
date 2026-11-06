package subscription

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	BillingHistoryUsecase interface {
		Execute(ctx context.Context, teacherID string) (*BillingHistoryOutput, apperrors.ApplicationError)
	}

	billingHistoryUsecase struct {
		contextFactory appcontext.Factory
	}

	BillingCycleData struct {
		SubscriptionID int     `json:"subscription_id"`
		Status         string  `json:"status"`
		Amount         float64 `json:"amount"`
		Currency       string  `json:"currency"`
		PeriodStart    string  `json:"period_start"`
		PeriodEnd      string  `json:"period_end"`
	}

	BillingHistoryData struct {
		AmountOwed *float64           `json:"amount_owed"`
		Currency   string             `json:"currency"`
		Cycles     []BillingCycleData `json:"cycles"`
	}

	BillingHistoryOutput struct {
		Data BillingHistoryData `json:"data"`
	}
)

func NewBillingHistoryUsecase(contextFactory appcontext.Factory) BillingHistoryUsecase {
	return &billingHistoryUsecase{contextFactory: contextFactory}
}

func (u *billingHistoryUsecase) Execute(ctx context.Context, teacherID string) (*BillingHistoryOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}

	cycles, err := app.Integrations.Payments.ListBillingCycles(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}

	data := BillingHistoryData{Cycles: make([]BillingCycleData, 0, len(cycles))}
	for _, cycle := range cycles {
		data.Cycles = append(data.Cycles, BillingCycleData{
			SubscriptionID: cycle.SubscriptionID,
			Status:         cycle.Status,
			Amount:         cycle.Amount,
			Currency:       cycle.Currency,
			PeriodStart:    cycle.PeriodStart.UTC().Format(time.RFC3339),
			PeriodEnd:      cycle.PeriodEnd.UTC().Format(time.RFC3339),
		})
	}
	if len(cycles) > 0 {
		data.Currency = cycles[0].Currency
	}
	if scope.GraceEndsAt != nil {
		for _, cycle := range cycles {
			if cycle.Status == "rejected" || cycle.Status == "cancelled" {
				amount := cycle.Amount
				data.AmountOwed = &amount
				data.Currency = cycle.Currency
				break
			}
		}
	}

	return &BillingHistoryOutput{Data: data}, nil
}

package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type billingPayments struct {
	payments.Client
	entitlement   *payments.Entitlement
	subscriptions []payments.Subscription
	plans         []payments.Plan
	cycles        []payments.BillingCycle
	cyclesErr     error
	userID        string
}

func billingHistoryApp(pay *billingPayments) *appcontext.Context {
	return &appcontext.Context{
		Repositories: &repositories.Repositories{School: subscriptionSchool(0), UserProfile: withinTrial()},
		Integrations: &integrations.Integrations{Payments: pay},
	}
}

func (f *billingPayments) GetEntitlement(context.Context, string) (*payments.Entitlement, error) {
	return f.entitlement, nil
}

func (f *billingPayments) ListSubscriptions(context.Context, string) ([]payments.Subscription, error) {
	return f.subscriptions, nil
}

func (f *billingPayments) ListPlans(context.Context) ([]payments.Plan, error) {
	return f.plans, nil
}

func (f *billingPayments) ListBillingCycles(_ context.Context, userID string) ([]payments.BillingCycle, error) {
	f.userID = userID
	return f.cycles, f.cyclesErr
}

func TestBillingHistoryUsesLatestFailedCycleOnlyInGrace(t *testing.T) {
	paidUntil := time.Now().UTC().Add(-24 * time.Hour)
	f := &billingPayments{
		subscriptions: []payments.Subscription{{PlanID: 1, CurrentPeriodEnd: &payments.Timestamp{Time: paidUntil}}},
		plans:         []payments.Plan{{ID: 1, Metadata: map[string]any{"max_students": float64(5)}}},
		cycles: []payments.BillingCycle{
			{SubscriptionID: 2, Status: "pending", Amount: 300, Currency: "ARS", PeriodStart: payments.Timestamp{Time: paidUntil}, PeriodEnd: payments.Timestamp{Time: paidUntil.AddDate(0, 1, 0)}},
			{SubscriptionID: 2, Status: "rejected", Amount: 250, Currency: "ARS", PeriodStart: payments.Timestamp{Time: paidUntil.AddDate(0, -1, 0)}, PeriodEnd: payments.Timestamp{Time: paidUntil}},
			{SubscriptionID: 1, Status: "cancelled", Amount: 100, Currency: "ARS", PeriodStart: payments.Timestamp{Time: paidUntil.AddDate(0, -2, 0)}, PeriodEnd: payments.Timestamp{Time: paidUntil.AddDate(0, -1, 0)}},
		},
	}
	app := billingHistoryApp(f)
	out, appErr := NewBillingHistoryUsecase(func() *appcontext.Context { return app }).Execute(context.Background(), "teacher-1")
	if appErr != nil {
		t.Fatalf("Execute: %v", appErr)
	}
	if f.userID != "teacher-1" || out.Data.AmountOwed == nil || *out.Data.AmountOwed != 250 {
		t.Fatalf("user = %q, owed = %v", f.userID, out.Data.AmountOwed)
	}
	if len(out.Data.Cycles) != 3 || out.Data.Cycles[1].PeriodStart == "" {
		t.Fatalf("cycles = %+v", out.Data.Cycles)
	}

	f.cycles = nil
	out, appErr = NewBillingHistoryUsecase(func() *appcontext.Context { return app }).Execute(context.Background(), "teacher-1")
	if appErr != nil || out.Data.AmountOwed != nil || out.Data.Cycles == nil {
		t.Fatalf("empty history: output=%+v error=%v", out, appErr)
	}

	f.entitlement = &payments.Entitlement{Active: true, Metadata: map[string]any{"max_students": float64(5)}}
	f.cycles = []payments.BillingCycle{{Status: "rejected", Amount: 250}}
	out, appErr = NewBillingHistoryUsecase(func() *appcontext.Context { return app }).Execute(context.Background(), "teacher-1")
	if appErr != nil || out.Data.AmountOwed != nil {
		t.Fatalf("outside grace: output=%+v error=%v", out, appErr)
	}
}

func TestBillingHistoryReportsPaymentServiceFailure(t *testing.T) {
	f := &billingPayments{entitlement: &payments.Entitlement{Active: true}, cyclesErr: errors.New("unavailable")}
	app := billingHistoryApp(f)
	_, appErr := NewBillingHistoryUsecase(func() *appcontext.Context { return app }).Execute(context.Background(), "teacher-1")
	if appErr == nil {
		t.Fatal("expected payment service error")
	}
}

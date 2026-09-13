package subscription

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

type plansUsecaseStub struct {
	output *ucSubscription.PlansOutput
}

func (s plansUsecaseStub) List(context.Context) (*ucSubscription.PlansOutput, apperrors.ApplicationError) {
	return s.output, nil
}

func (plansUsecaseStub) Create(context.Context, ucSubscription.PlanInput) (*ucSubscription.PlanOutput, apperrors.ApplicationError) {
	return nil, nil
}

func (plansUsecaseStub) Update(context.Context, int, ucSubscription.PlanInput) (*ucSubscription.PlanOutput, apperrors.ApplicationError) {
	return nil, nil
}

func (plansUsecaseStub) Deactivate(context.Context, int) (*ucSubscription.PlanOutput, apperrors.ApplicationError) {
	return nil, nil
}

func TestPublicListPlansHidesRetiredPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := gin.New()
	app.GET("/api/public/subscription-plans", NewPublicListPlansHandler(plansUsecaseStub{
		output: &ucSubscription.PlansOutput{Data: []ucSubscription.CatalogPlanData{
			{PlanID: 1, Name: "Disponible", Active: true},
			{PlanID: 2, Name: "Retirado", Active: false},
		}},
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/public/subscription-plans", nil)
	res := httptest.NewRecorder()
	app.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got, want := res.Body.String(), `{"data":[{"plan_id":1,"name":"Disponible","amount":0,"currency":"","interval":"","max_students":0,"active":true}]}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

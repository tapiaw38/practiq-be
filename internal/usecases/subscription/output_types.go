package subscription

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

type (
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
)

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

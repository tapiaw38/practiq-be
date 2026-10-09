package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewPublicListPlansHandler(uc ucSubscription.ListPlansUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		active := make([]ucSubscription.CatalogPlanData, 0, len(output.Data))
		for _, plan := range output.Data {
			if plan.Active {
				active = append(active, plan)
			}
		}
		c.JSON(http.StatusOK, publicListPlansOutput{Data: active})
	}
}

type publicListPlansOutput struct {
	Data []ucSubscription.CatalogPlanData `json:"data"`
}

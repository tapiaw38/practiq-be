package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewCreatePlanHandler(uc ucSubscription.CreatePlanUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubscription.CreatePlanInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}

package subscription

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewUpdatePlanHandler(uc ucSubscription.UpdatePlanUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		planID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid plan id"})
			return
		}
		var input ucSubscription.UpdatePlanInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Execute(c, planID, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

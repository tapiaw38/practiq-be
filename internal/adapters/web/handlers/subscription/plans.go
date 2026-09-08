package subscription

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewListPlansHandler(uc ucSubscription.PlansUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.List(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewCreatePlanHandler(uc ucSubscription.PlansUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubscription.PlanInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Create(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}

func NewUpdatePlanHandler(uc ucSubscription.PlansUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		planID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid plan id"})
			return
		}
		var input ucSubscription.PlanInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Update(c, planID, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewDeactivatePlanHandler(uc ucSubscription.PlansUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		planID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid plan id"})
			return
		}
		output, appErr := uc.Deactivate(c, planID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

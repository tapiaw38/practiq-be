package learningstrategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
)

type assignToCourseInput struct {
	StrategyID string `json:"strategy_id" binding:"required"`
	IsDefault  bool   `json:"is_default"`
	Config     string `json:"config"`
}

func NewAssignToCourseHandler(uc ucLS.AssignToCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		courseID := c.Param("id")
		isAdmin := middlewares.HasRole(c, "admin") || middlewares.HasRole(c, "superadmin")

		var input assignToCourseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, courseID, isAdmin, ucLS.AssignToCourseInput{
			StrategyID: input.StrategyID,
			IsDefault:  input.IsDefault,
			Config:     input.Config,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

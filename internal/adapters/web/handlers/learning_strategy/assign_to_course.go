package learningstrategy

import (
	"net/http"

	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"

	"github.com/gin-gonic/gin"
)

type assignToCourseInput struct {
	StrategyID string `json:"strategy_id" binding:"required"`
	IsDefault  bool   `json:"is_default"`
	Config     string `json:"config"`
}

func NewAssignToCourseHandler(uc ucLS.AssignToCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input assignToCourseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, courseID, ucLS.AssignToCourseInput{
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

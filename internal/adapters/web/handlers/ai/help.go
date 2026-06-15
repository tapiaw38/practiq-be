package ai

import (
	"net/http"

	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type helpInput struct {
	ExerciseID     string `json:"exercise_id"`
	Question       string `json:"question" binding:"required"`
	HelpType       string `json:"help_type"`
	ConversationID string `json:"conversation_id"`
}

func NewHelpHandler(uc ucAI.HelpUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		var input helpInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucAI.HelpInput{
			StudentID:      studentID,
			ExerciseID:     input.ExerciseID,
			Question:       input.Question,
			HelpType:       input.HelpType,
			ConversationID: input.ConversationID,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

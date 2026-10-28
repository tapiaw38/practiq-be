package ai

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

type copilotInput struct {
	ExerciseID    string `json:"exercise_id"`
	ContextID     string `json:"context_id"`
	Question      string `json:"question" binding:"required"`
	StudentAnswer string `json:"student_answer"`
	Intent        string `json:"intent"`
}

func NewCopilotHandler(help ucAI.HelpUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input copilotInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		intent := strings.ToLower(strings.TrimSpace(input.Intent))
		if intent != "hint" && intent != "explanation" && intent != "similar_example" && intent != "review_answer" {
			intent = "hint"
		}
		output, appErr := help.Execute(c, middlewares.GetUserID(c), ucAI.HelpInput{ExerciseID: input.ExerciseID, Question: input.Question, StudentAnswer: input.StudentAnswer, HelpType: intent})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"context_id": input.ContextID,
			"blocks":     []gin.H{{"type": intent, "content": output.Data.Response}},

			"suggested_actions": []gin.H{
				{"id": "hint", "type": "prompt", "label": "Otra pista", "prompt": "Dame otra pista sin revelar la respuesta."},
				{"id": "explanation", "type": "prompt", "label": "Explicame", "prompt": "Explicame paso a paso usando el ejercicio actual."},
				{"id": "similar_example", "type": "prompt", "label": "Ejemplo", "prompt": "Dame un ejemplo similar con números diferentes."},
				{"id": "review_answer", "type": "prompt", "label": "Revisá", "prompt": "Revisá mi respuesta actual."},
			},
		}})
	}
}

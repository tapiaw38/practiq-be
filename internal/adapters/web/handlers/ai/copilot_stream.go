package ai

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func NewCopilotStreamHandler(help ucAI.HelpUsecase) gin.HandlerFunc {
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

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.SSEvent("status", gin.H{"message": "Ana está preparando ayuda para este ejercicio…"})
		c.Writer.Flush()

		output, appErr := help.Execute(c, middlewares.GetUserID(c), ucAI.HelpInput{ExerciseID: input.ExerciseID, Question: input.Question, StudentAnswer: input.StudentAnswer, HelpType: intent})
		if appErr != nil {
			appErr.Log(c)
			c.SSEvent("error", gin.H{"message": "No se pudo obtener ayuda ahora."})
			c.Writer.Flush()
			return
		}
		c.SSEvent("response", gin.H{"data": gin.H{
			"context_id": input.ContextID,
			"blocks":     []gin.H{{"type": intent, "content": output.Data.Response}},
			"suggested_actions": []gin.H{
				{"id": "hint", "type": "prompt", "label": "Otra pista", "prompt": "Dame otra pista sin revelar la respuesta."},
				{"id": "explanation", "type": "prompt", "label": "Explicame", "prompt": "Explicame paso a paso usando el ejercicio actual."},
				{"id": "similar_example", "type": "prompt", "label": "Ejemplo", "prompt": "Dame un ejemplo similar con números diferentes."},
				{"id": "review_answer", "type": "prompt", "label": "Revisá", "prompt": "Revisá mi respuesta actual."},
			},
		}})
		c.Writer.Flush()
	}
}

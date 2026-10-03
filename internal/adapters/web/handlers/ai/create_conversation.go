package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func NewCreateConversationHandler(uc ucAI.CreateConversationUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		var input ucAI.CreateConversationInput
		c.ShouldBindJSON(&input)

		output, appErr := uc.Execute(c, studentID, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

package ai

import (
	"net/http"

	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"

	"github.com/gin-gonic/gin"
)

func NewGetMessagesHandler(uc ucAI.GetMessagesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		output, appErr := uc.Execute(c, conversationID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

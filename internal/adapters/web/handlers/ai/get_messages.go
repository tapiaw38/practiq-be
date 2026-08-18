package ai

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"

	"github.com/gin-gonic/gin"
)

func NewGetMessagesHandler(uc ucAI.GetMessagesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		output, appErr := uc.Execute(c, conversationID, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

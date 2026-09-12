package topic

import (
	"net/http"

	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type updateTopicInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
}

func NewUpdateHandler(uc ucTopic.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		var input updateTopicInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, id, ucTopic.UpdateInput{
			Title:       input.Title,
			Description: input.Description,
			OrderIndex:  input.OrderIndex,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

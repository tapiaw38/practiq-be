package topic

import (
	"net/http"

	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucTopic.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

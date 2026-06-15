package exercise

import (
	"net/http"

	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucExercise.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("id")
		output, appErr := uc.Execute(c, topicID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

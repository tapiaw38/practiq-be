package learningstrategy

import (
	"net/http"

	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"

	"github.com/gin-gonic/gin"
)

func NewListByCourseHandler(uc ucLS.ListByCourseUsecase) gin.HandlerFunc {
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

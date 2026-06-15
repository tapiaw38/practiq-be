package learningstrategy

import (
	"net/http"

	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"

	"github.com/gin-gonic/gin"
)

func NewUnassignFromCourseHandler(uc ucLS.UnassignFromCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "strategy unassigned from course"})
	}
}

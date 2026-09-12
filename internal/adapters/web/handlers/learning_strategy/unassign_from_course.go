package learningstrategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
)

func NewUnassignFromCourseHandler(uc ucLS.UnassignFromCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		id := c.Param("id")
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		if appErr := uc.Execute(c, requesterID, id, isSuperAdmin); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "strategy unassigned from course"})
	}
}

package studentprogress

import (
	"net/http"

	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewGetCourseProgressHandler(uc ucProgress.GetCourseProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, studentID, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

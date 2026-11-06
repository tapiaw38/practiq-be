package courseprogress

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucCP "github.com/tapiaw38/practiq-be/internal/usecases/course_progress"
)

func NewGetForStudentHandler(uc ucCP.GetForStudentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		studentID := c.Param("studentId")
		courseID := c.Param("courseId")
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		output, appErr := uc.Execute(c, requesterID, studentID, courseID, isSuperAdmin)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

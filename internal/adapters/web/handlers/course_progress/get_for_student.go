package courseprogress

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucCP "github.com/tapiaw38/practiq-be/internal/usecases/course_progress"
)

func NewGetForStudentHandler(uc ucCP.GetForStudentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		courseID := c.Param("courseId")

		output, appErr := uc.Execute(c, studentID, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

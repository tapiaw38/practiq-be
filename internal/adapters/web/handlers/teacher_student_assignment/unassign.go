package teacherstudentassignment

import (
	"net/http"

	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"

	"github.com/gin-gonic/gin"
)

func NewUnassignHandler(uc ucAssignment.UnassignUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := c.Param("teacherId")
		studentID := c.Param("studentId")
		output, appErr := uc.Execute(c, teacherID, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

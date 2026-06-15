package teacherstudentassignment

import (
	"net/http"

	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"

	"github.com/gin-gonic/gin"
)

func NewListStudentsHandler(uc ucAssignment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := c.Param("teacherId")
		output, appErr := uc.Execute(c, teacherID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

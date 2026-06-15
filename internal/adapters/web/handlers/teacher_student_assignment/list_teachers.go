package teacherstudentassignment

import (
	"net/http"

	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"

	"github.com/gin-gonic/gin"
)

func NewListTeachersHandler(uc ucAssignment.ListTeachersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		output, appErr := uc.Execute(c, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

package teacherstudentassignment

import (
	"net/http"

	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewListMyStudentsHandler(uc ucAssignment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := middlewares.GetUserID(c)

		input := ucAssignment.ListStudentsInput{
			TeacherID: teacherID,
		}

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

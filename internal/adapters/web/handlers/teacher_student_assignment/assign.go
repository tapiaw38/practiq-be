package teacherstudentassignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"
)

type assignInput struct {
	TeacherID string `json:"teacher_id" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}

func NewAssignHandler(uc ucAssignment.AssignUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input assignInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, input.TeacherID, input.StudentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

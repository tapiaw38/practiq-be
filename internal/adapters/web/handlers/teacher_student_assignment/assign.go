package teacherstudentassignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"
)

func NewAssignHandler(uc ucAssignment.AssignUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucAssignment.AssignInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

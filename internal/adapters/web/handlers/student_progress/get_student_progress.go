package studentprogress

import (
	"net/http"

	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"

	"github.com/gin-gonic/gin"
)

func NewGetStudentProgressHandler(uc ucProgress.GetStudentProgressUsecase) gin.HandlerFunc {
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

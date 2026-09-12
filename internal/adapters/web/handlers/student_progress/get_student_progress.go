package studentprogress

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"

	"github.com/gin-gonic/gin"
)

func NewGetStudentProgressHandler(uc ucProgress.GetStudentProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		studentID := c.Param("studentId")

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

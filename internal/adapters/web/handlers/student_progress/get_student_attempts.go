package studentprogress

import (
	"net/http"

	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"

	"github.com/gin-gonic/gin"
)

func NewGetStudentAttemptsHandler(uc ucProgress.GetStudentAttemptsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		sheetID := c.Query("sheet_id")
		if sheetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "sheet_id query param is required"})
			return
		}
		output, appErr := uc.Execute(c, studentID, sheetID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

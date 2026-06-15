package courseprogress

import (
	"net/http"

	ucCP "github.com/tapiaw38/practiq-be/internal/usecases/course_progress"

	"github.com/gin-gonic/gin"
)

func NewListForStudentHandler(uc ucCP.ListForStudentUsecase) gin.HandlerFunc {
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

package enrollment

import (
	"net/http"

	ucEnrollment "github.com/tapiaw38/practiq-be/internal/usecases/enrollment"

	"github.com/gin-gonic/gin"
)

func NewListStudentsHandler(uc ucEnrollment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

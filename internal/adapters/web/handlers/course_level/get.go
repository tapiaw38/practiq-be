package courselevel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucLevel "github.com/tapiaw38/practiq-be/internal/usecases/course_level"
)

func NewGetHandler(uc ucLevel.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		studentID := middlewares.GetUserID(c)

		output, appErr := uc.Execute(c, courseID, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

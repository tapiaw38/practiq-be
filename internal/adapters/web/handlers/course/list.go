package course

import (
	"net/http"

	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewListHandler(uc ucCourse.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)
		role := c.Query("role")

		input := ucCourse.ListInput{}
		if role == "teacher" {
			input.TeacherID = userID
		} else {
			input.StudentID = userID
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

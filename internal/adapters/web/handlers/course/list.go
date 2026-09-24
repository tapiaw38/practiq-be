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
			// A platform superadmin opening a school is operating it, not
			// teaching in it. Narrowing to the courses they own left them
			// looking at an empty school and re-creating by hand what was
			// already there. The school header still bounds the answer, the
			// same way it already does for grades and subjects.
			if !middlewares.IsSuperAdmin(c) {
				input.TeacherID = userID
			}
		} else {
			input.StudentID = userID
		}
		input.SchoolID = c.GetHeader("X-School-ID")

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

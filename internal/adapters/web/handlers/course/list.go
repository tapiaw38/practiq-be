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

		teacherID, studentID := "", ""
		if role == "teacher" {

			if !middlewares.IsSuperAdmin(c) {
				teacherID = userID
			}
		} else {
			studentID = userID
		}

		output, appErr := uc.Execute(c, teacherID, studentID, c.GetHeader("X-School-ID"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

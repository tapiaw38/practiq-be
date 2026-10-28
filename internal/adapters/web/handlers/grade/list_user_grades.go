package grade

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

func NewListUserGradesHandler(uc ucGrade.ListUserGradesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")

		if userID != middlewares.GetUserID(c) && !middlewares.IsTeacher(c) {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "forbidden"})
			return
		}

		output, appErr := uc.Execute(c, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

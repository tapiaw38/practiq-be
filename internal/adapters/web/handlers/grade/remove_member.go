package grade

import (
	"net/http"

	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

func NewRemoveMemberHandler(uc ucGrade.RemoveMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		userID := c.Param("userId")
		output, appErr := uc.Execute(c, gradeID, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

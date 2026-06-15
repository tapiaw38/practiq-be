package grade

import (
	"net/http"

	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

func NewListUserGradesHandler(uc ucGrade.ListUserGradesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")
		output, appErr := uc.Execute(c, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

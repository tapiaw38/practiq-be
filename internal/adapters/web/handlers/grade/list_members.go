package grade

import (
	"net/http"

	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

func NewListMembersHandler(uc ucGrade.ListMembersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		output, appErr := uc.Execute(c, gradeID, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

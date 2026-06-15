package practicesheet

import (
	"net/http"

	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucPS.ListUsecase) gin.HandlerFunc {
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

package school

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func NewListHandler(uc ucSchool.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

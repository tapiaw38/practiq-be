package grade

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucGrade.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.GetHeader("X-School-ID"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

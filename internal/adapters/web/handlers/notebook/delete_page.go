package notebook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

func NewDeletePageHandler(uc ucNB.DeletePageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

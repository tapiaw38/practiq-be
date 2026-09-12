package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewListHandler(uc ucNB.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		out, appErr := uc.Execute(c, requesterID, isSuperAdmin, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

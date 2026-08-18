package practicesheet

import (
	"net/http"

	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewGetHandler(uc ucPS.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

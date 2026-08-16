package practicesheet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
)

func NewGetAssistantMediaHandler(uc ucPS.GetAssistantMediaUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(
			c,
			middlewares.GetUserID(c),
			middlewares.HasRole(c, "admin", "superadmin"),
			c.Param("id"),
			c.Param("exerciseId"),
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Header("Cache-Control", "private, max-age=300")
		c.Data(http.StatusOK, output.ContentType, output.Content)
	}
}

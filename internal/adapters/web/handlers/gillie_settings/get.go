package gilliesettings

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucGillie "github.com/tapiaw38/practiq-be/internal/usecases/gillie_settings"
)

func NewGetHandler(uc ucGillie.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

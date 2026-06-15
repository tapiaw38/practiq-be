package userprofile

import (
	"net/http"

	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewGetHandler(uc ucProfile.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

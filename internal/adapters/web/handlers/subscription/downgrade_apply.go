package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewDowngradeApplyHandler(uc ucSubscription.DowngradeApplyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body ucSubscription.DowngradeApplyInput

		_ = c.ShouldBindJSON(&body)

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), body)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

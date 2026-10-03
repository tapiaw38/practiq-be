package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewHostedCheckoutHandler(uc ucSubscription.HostedCheckoutUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubscription.HostedCheckoutInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		out, appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("Authorization"), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

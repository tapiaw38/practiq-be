package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewSubscribeHandler(uc ucSubscription.SubscribeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubscription.SubscribeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		// The teacher comes from the token, never from the body.
		if appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("Authorization"), input); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// NewCheckoutConfigHandler serves the gateway's public key, which the browser
// needs to turn a card into a token without the card passing through us.
func NewCheckoutConfigHandler(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ucSubscription.CheckoutConfig(publicKey))
	}
}

package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewCheckoutConfigHandler(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ucSubscription.CheckoutConfig(publicKey))
	}
}

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

		if appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("Authorization"), input); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

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

func NewChangePlanHandler(uc ucSubscription.ChangePlanUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubscription.ChangePlanInput
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

func NewCheckoutConfigHandler(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ucSubscription.CheckoutConfig(publicKey))
	}
}

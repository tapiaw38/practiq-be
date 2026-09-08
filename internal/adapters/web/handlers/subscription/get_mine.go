package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewGetMineHandler(uc ucSubscription.GetMineUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// The teacher is the authenticated caller, never a parameter: a
		// subscription is the one thing a user must not be able to read for
		// somebody else by guessing an id.
		output, appErr := uc.Execute(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

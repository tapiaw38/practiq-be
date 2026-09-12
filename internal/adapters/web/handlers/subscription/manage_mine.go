package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

// NewManageMineHandler binds one action to one route, so the action is decided
// by which endpoint was called and never by a value in the request.
func NewManageMineHandler(uc ucSubscription.ManageMineUsecase, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, middlewares.GetUserID(c), action); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

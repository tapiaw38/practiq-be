package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewListPlansHandler(uc ucSubscription.ListPlansUsecase) gin.HandlerFunc {
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

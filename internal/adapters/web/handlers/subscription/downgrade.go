package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func NewDowngradePreviewHandler(uc ucSubscription.DowngradeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Preview(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewDowngradeApplyHandler(uc ucSubscription.DowngradeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			// Keep is the teacher's choice of who stays. Empty means the
			// automatic order: least recently active go first.
			Keep []string `json:"keep"`
		}
		// A body is optional: applying without one takes the automatic order.
		_ = c.ShouldBindJSON(&body)

		output, appErr := uc.Apply(c, middlewares.GetUserID(c), body.Keep)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewReactivateStudentHandler(uc ucSubscription.DowngradeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Reactivate(c, middlewares.GetUserID(c), c.Param("studentId")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

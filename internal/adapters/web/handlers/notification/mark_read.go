package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNotification "github.com/tapiaw38/practiq-be/internal/usecases/notification"
)

func NewMarkReadHandler(uc ucNotification.MarkReadUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, c.Param("id"), middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewMarkAllReadHandler(uc ucNotification.MarkAllReadUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucNotification.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, c.Param("id"), middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

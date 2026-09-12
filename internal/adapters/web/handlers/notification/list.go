package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNotification "github.com/tapiaw38/practiq-be/internal/usecases/notification"
)

func NewListHandler(uc ucNotification.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))

		output, appErr := uc.Execute(c, ucNotification.ListInput{
			UserID:     middlewares.GetUserID(c),
			UnreadOnly: c.Query("unread_only") == "true",
			Limit:      limit,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

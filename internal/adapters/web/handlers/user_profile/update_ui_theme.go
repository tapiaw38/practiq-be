package userprofile

import (
	"net/http"

	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewUpdateUIThemeHandler(uc ucProfile.UpdateUIThemeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input uiThemeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, middlewares.IsSuperAdmin(c), ucProfile.UpdateUIThemeInput{
			ID:          userID,
			UITheme:     input.UITheme,
			BearerToken: c.GetHeader("Authorization"),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

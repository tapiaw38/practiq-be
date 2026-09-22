package userprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"
)

type avatarSeedInput struct {
	AvatarSeed string `json:"avatar_seed"`
}

// Self-service only: the id comes from the token, never from the body, so
// nobody repaints somebody else's avatar.
func NewUpdateAvatarSeedHandler(uc ucProfile.UpdateAvatarSeedUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input avatarSeedInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, middlewares.IsSuperAdmin(c), ucProfile.UpdateAvatarSeedInput{
			ID:          userID,
			AvatarSeed:  input.AvatarSeed,
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

package userprofile

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
)

func NewGetByIDHandler(uc ucProfile.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		profileID := c.Param("id")

		if profileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "profile id required"})
			return
		}

		// Authorization check: only admins or the profile owner can view profiles
		if !isAdmin && requesterID != profileID {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "cannot view other user profiles"})
			return
		}

		output, appErr := uc.Execute(c, profileID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

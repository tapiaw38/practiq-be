package impersonation

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	profileRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/user_profile"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/platform/auth"
)

type startInput struct {
	UserID string `json:"user_id"`
}

// Start is mounted only under superadmin routes. Session expires in 15 minutes
// and has no refresh token or target roles.
func Start(profiles profileRepo.Repository, authAPI authapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input startInput
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.UserID) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "user_id required"})
			return
		}
		targetID, operatorID := strings.TrimSpace(input.UserID), middlewares.GetUserID(c)
		if targetID == operatorID {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "cannot impersonate current user"})
			return
		}
		profile, err := profiles.Get(c.Request.Context(), targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "failed to load target profile"})
			return
		}
		if profile == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "common:not-found", "message": "user has not initialized Practiq yet"})
			return
		}
		version, err := authAPI.GetTokenVersion(c.Request.Context(), c.GetHeader("Authorization"), targetID)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"code": "impersonation:auth-unavailable", "message": "could not verify target session"})
			return
		}
		token, err := auth.GenerateReadOnlyImpersonationToken(targetID, version, operatorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "could not create impersonation session"})
			return
		}
		log.Printf("[impersonation] operator=%q target=%q read_only=true", operatorID, targetID)
		c.JSON(http.StatusOK, gin.H{"token": token, "user_id": targetID, "profile_type": profile.ProfileType, "read_only": true})
	}
}

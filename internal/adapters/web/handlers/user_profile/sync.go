package userprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"
)

type syncInput struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	ProfileType string `json:"profile_type"`
	// Timezone is the IANA zone the browser reports. Empty leaves whatever is
	// stored: the streak is measured with it and most clients do not send one.
	Timezone         string `json:"timezone"`
	AssistantBaseURL string `json:"assistant_base_url"`
	AssistantAPIKey  string `json:"assistant_api_key"`
}
type assistantConfigInput struct {
	AssistantBaseURL string `json:"assistant_base_url"`
	AssistantAPIKey  string `json:"assistant_api_key"`
}
type academicStatusInput struct {
	AcademicStatus string `json:"academic_status" binding:"required"`
}

func NewSyncHandler(uc ucProfile.SyncUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input syncInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, ucProfile.SyncInput{
			ID:               userID,
			Name:             input.Name,
			Email:            input.Email,
			ProfileType:      input.ProfileType,
			Timezone:         input.Timezone,
			AssistantBaseURL: input.AssistantBaseURL,
			AssistantAPIKey:  input.AssistantAPIKey,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

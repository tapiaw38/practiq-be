package userprofile

import (
	"net/http"

	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
)

func NewUpdateAssistantConfigByIDHandler(uc ucProfile.UpdateAssistantConfigUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input assistantConfigInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		profileID := c.Param("id")
		if profileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "profile id required"})
			return
		}

		output, appErr := uc.Execute(c, ucProfile.UpdateAssistantConfigInput{
			ID:               profileID,
			AssistantBaseURL: input.AssistantBaseURL,
			AssistantAPIKey:  input.AssistantAPIKey,
			UITheme:          input.UITheme,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

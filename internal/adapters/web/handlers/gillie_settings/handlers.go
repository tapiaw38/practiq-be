package gilliesettings

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	repo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/gillie_settings"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/platform/secretbox"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
)

type input struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

func Get(r repo.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		stored, err := r.Get(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "gillie:read-error", "message": "could not read the assistant settings"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"base_url":   stored.BaseURL,
			"api_key":    masked(stored.APIKeyLast4),
			"configured": stored.BaseURL != "" && stored.APIKeyEncrypted != "",
			"updated_by": stored.UpdatedBy,
		}})
	}
}

func Update(r repo.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in input
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		baseURL := strings.TrimRight(strings.TrimSpace(in.BaseURL), "/")
		if baseURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "gillie:bad-url", "message": "base_url is required"})
			return
		}

		if err := assistant.ValidateBaseURL(baseURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "gillie:bad-url", "message": err.Error()})
			return
		}

		settings := repo.Settings{BaseURL: baseURL, UpdatedBy: middlewares.GetUserID(c)}

		if apiKey := strings.TrimSpace(in.APIKey); apiKey != "" {
			sealed, err := assistantcfg.Seal(apiKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "gillie:encryption-unavailable",
					"message": "the server cannot encrypt the key: set GILLIE_CONFIG_SECRET",
				})
				return
			}
			settings.APIKeyEncrypted = sealed
			settings.APIKeyLast4 = secretbox.Last4(apiKey)
		}

		if err := r.Save(c, settings); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "gillie:save-error", "message": "could not save the assistant settings"})
			return
		}

		stored, err := r.Get(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "gillie:read-error", "message": "could not read the assistant settings"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"base_url":   stored.BaseURL,
			"api_key":    masked(stored.APIKeyLast4),
			"configured": stored.BaseURL != "" && stored.APIKeyEncrypted != "",
			"updated_by": stored.UpdatedBy,
		}})
	}
}

func masked(last4 string) string {
	if last4 == "" {
		return ""
	}
	return "••••••••" + last4
}

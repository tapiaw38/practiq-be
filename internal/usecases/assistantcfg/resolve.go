package assistantcfg

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
	"github.com/tapiaw38/practiq-be/internal/platform/secretbox"
)

var (
	boxOnce sync.Once
	box     *secretbox.Box
)

func encryptionBox() *secretbox.Box {
	boxOnce.Do(func() {
		cfg := config.GetConfigService()
		opened, err := secretbox.New(cfg.ServerConfig.GillieConfigSecret)
		if err != nil {
			log.Printf("[assistant] GILLIE_CONFIG_SECRET unusable: %v", err)
			return
		}
		box = opened
	})
	return box
}

// Seal encrypts a Gillie API key for storage.
func Seal(apiKey string) (string, error) {
	b := encryptionBox()
	if b == nil {
		return "", secretbox.ErrNoKey
	}
	return b.Seal(apiKey)
}

// Resolve returns the platform's Gillie configuration, which a superadmin
// stores and nothing else supplies. An empty Config reads as "not configured"
// everywhere it is used, so every failure here turns the assistant off rather
// than reaching for a second set of credentials nobody chose.
func Resolve(ctx context.Context, app *appcontext.Context) assistant.Config {
	if app == nil || app.Repositories == nil || app.Repositories.GillieSettings == nil {
		return assistant.Config{}
	}

	stored, err := app.Repositories.GillieSettings.Get(ctx)
	if err != nil {
		log.Printf("[assistant] could not read gillie settings: %v", err)
		return assistant.Config{}
	}
	if stored.APIKeyEncrypted == "" {
		return assistant.Config{BaseURL: strings.TrimSpace(stored.BaseURL)}
	}

	b := encryptionBox()
	if b == nil {
		log.Print("[assistant] GILLIE_CONFIG_SECRET is unusable: the stored API key cannot be read")
		return assistant.Config{}
	}
	apiKey, err := b.Open(stored.APIKeyEncrypted)
	if err != nil {
		log.Printf("[assistant] the stored API key cannot be decrypted, set it again from the admin panel: %v", err)
		return assistant.Config{}
	}

	return assistant.Config{
		BaseURL: strings.TrimSpace(stored.BaseURL),
		APIKey:  apiKey,
	}
}

// Enabled reports whether the platform assistant can be used, without handing
// out the credentials that answer the question.
func Enabled(ctx context.Context, app *appcontext.Context) bool {
	cfg := Resolve(ctx, app)
	return strings.TrimSpace(cfg.BaseURL) != "" && strings.TrimSpace(cfg.APIKey) != ""
}

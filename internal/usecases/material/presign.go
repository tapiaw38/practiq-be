package material

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

const viewLinkTTL = time.Hour

const materialsFolder = "materials"

func withViewURL(app *appcontext.Context, data MaterialData) MaterialData {
	if data.FileURL == "" || app.ImageStorage == nil {
		return data
	}
	if signed, ok := app.ImageStorage.PresignGetURL(data.FileURL, viewLinkTTL); ok {
		data.ViewURL = signed
	}
	return data
}

package material

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// viewLinkTTL has to outlive a reading session without leaving a shareable
// link around for long. Reloading the course view issues fresh ones.
const viewLinkTTL = time.Hour

// withViewURL fills the temporary URL the browser can actually open. FileURL
// stays canonical so it can be written back unchanged.
func withViewURL(app *appcontext.Context, data MaterialData) MaterialData {
	if data.FileURL == "" || app.ImageStorage == nil {
		return data
	}
	if signed, ok := app.ImageStorage.PresignGetURL(data.FileURL, viewLinkTTL); ok {
		data.ViewURL = signed
	}
	return data
}

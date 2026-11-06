package datauri

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func Resolve(ctx context.Context, app *appcontext.Context, value string) string {
	if app == nil || app.ImageStorage == nil {
		return value
	}
	resolved, err := app.ImageStorage.ResolveDataURI(ctx, value)
	if err != nil {
		return value
	}
	return resolved
}

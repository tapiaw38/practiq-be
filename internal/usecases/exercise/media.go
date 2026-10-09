package exercise

import (
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
)

func validateExerciseMediaURL(app *appcontext.Context, requesterID, metadata, previousMetadata string) apperrors.ApplicationError {
	mediaURL := (domain.Exercise{Metadata: metadata}).MediaURL()
	if mediaURL == "" || mediaURL == (domain.Exercise{Metadata: previousMetadata}).MediaURL() {
		return nil
	}
	if app == nil || app.ImageStorage == nil || !app.ImageStorage.OwnsFileURL(mediaURL, "exercises", requesterID) {
		return apperrors.NewBadRequestError("exercise media must be uploaded by the teacher editing this exercise")
	}
	return nil
}

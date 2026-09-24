package exercise

import (
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
)

// validateExerciseMediaURL only permits a newly attached medium when it was
// uploaded by the teacher editing the exercise. Existing media is allowed on
// updates so an administrator can edit legacy/shared exercises safely.
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

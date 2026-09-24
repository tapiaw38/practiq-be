package studentprogress

import (
	"context"
	"log"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// studentLocation resolves the zone a student's streak is measured in.
//
// A missing profile or an unreadable one is not worth failing the request over:
// the progress is still correct, only the day boundary falls back to the
// product default.
func studentLocation(ctx context.Context, app *appcontext.Context, studentID string) *time.Location {
	profile, err := app.Repositories.UserProfile.Get(ctx, studentID)
	if err != nil {
		log.Printf("[student_progress] could not read timezone student_id=%s err=%v", studentID, err)
		return domain.StudentLocation("")
	}
	if profile == nil {
		return domain.StudentLocation("")
	}
	return domain.StudentLocation(profile.Timezone)
}

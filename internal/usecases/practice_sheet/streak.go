package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func currentStudentStreak(ctx context.Context, app *appcontext.Context, studentID string) (int, error) {
	profile, err := app.Repositories.UserProfile.Get(ctx, studentID)
	if err != nil {
		return 0, err
	}
	location := domain.StudentLocation("")
	if profile != nil {
		location = domain.StudentLocation(profile.Timezone)
	}
	progress, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return 0, err
	}
	return domain.CurrentStreak(progress, location), nil
}

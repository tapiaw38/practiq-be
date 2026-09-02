package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	// DashboardUsecase serves the student home in one call.
	//
	// The screen used to assemble itself from about eighteen requests five
	// round trips deep: profile, then courses and progress, then practice
	// sheets, notebooks and levels once per course. The server spent about
	// 70ms of that; the rest was latency between Argentina and the API.
	//
	// It is named for the screen rather than something like "statistics" on
	// purpose. A view-shaped endpoint can change with its view; a generic one
	// accumulates half-used fields from every caller and stops being safe to
	// touch.
	DashboardUsecase interface {
		Execute(ctx context.Context, studentID string) (*DashboardOutput, apperrors.ApplicationError)
	}

	dashboardUsecase struct {
		contextFactory appcontext.Factory
	}

	CourseSummaryData struct {
		CourseID       string   `json:"course_id"`
		Title          string   `json:"title"`
		Subject        string   `json:"subject"`
		PracticeSheets int      `json:"practice_sheets"`
		LevelTests     int      `json:"level_tests"`
		Notebooks      int      `json:"notebooks"`
		CurrentLevel   int      `json:"current_level"`
		TopicIDs       []string `json:"topic_ids"`
	}

	DashboardData struct {
		Courses  []CourseSummaryData `json:"courses"`
		Progress []ProgressData      `json:"progress"`
		// StreakDays is the student's best live streak. It goes through the
		// domain rule rather than a SQL MAX so a streak the student already
		// broke is not reported.
		StreakDays           int    `json:"streak_days"`
		LastPracticedSheetID string `json:"last_practiced_sheet_id,omitempty"`
	}

	DashboardOutput struct {
		Data DashboardData `json:"data"`
	}
)

func NewDashboardUsecase(contextFactory appcontext.Factory) DashboardUsecase {
	return &dashboardUsecase{contextFactory: contextFactory}
}

func (u *dashboardUsecase) Execute(ctx context.Context, studentID string) (*DashboardOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	summaries, err := app.Repositories.Course.ListDashboardSummaries(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	progress, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	lastSheetID, err := app.Repositories.StudentAttempt.GetLastPracticedSheetID(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptGetError, err)
	}

	loc := studentLocation(ctx, app, studentID)

	courses := make([]CourseSummaryData, 0, len(summaries))
	for _, s := range summaries {
		courses = append(courses, CourseSummaryData{
			CourseID:       s.CourseID,
			Title:          s.Title,
			Subject:        s.Subject,
			PracticeSheets: s.PracticeSheets,
			LevelTests:     s.LevelTests,
			Notebooks:      s.Notebooks,
			CurrentLevel:   s.CurrentLevel,
			TopicIDs:       s.TopicIDs,
		})
	}

	progressData := make([]ProgressData, 0, len(progress))
	streak := 0
	for _, p := range progress {
		progressData = append(progressData, toProgressData(p, loc))
		if s := domain.EffectiveStreak(p, loc); s > streak {
			streak = s
		}
	}

	return &DashboardOutput{Data: DashboardData{
		Courses:              courses,
		Progress:             progressData,
		StreakDays:           streak,
		LastPracticedSheetID: lastSheetID,
	}}, nil
}

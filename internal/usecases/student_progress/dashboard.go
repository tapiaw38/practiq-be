package studentprogress

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DashboardUsecase interface {
		Execute(ctx context.Context, studentID string) (*DashboardOutput, apperrors.ApplicationError)
	}

	dashboardUsecase struct {
		contextFactory appcontext.Factory
	}

	CourseSummaryData struct {
		CourseID string `json:"course_id"`

		SchoolID       string   `json:"school_id,omitempty"`
		SchoolName     string   `json:"school_name,omitempty"`
		Title          string   `json:"title"`
		Subject        string   `json:"subject"`
		GradeName      string   `json:"grade_name,omitempty"`
		PracticeSheets int      `json:"practice_sheets"`
		LevelTests     int      `json:"level_tests"`
		Notebooks      int      `json:"notebooks"`
		CurrentLevel   int      `json:"current_level"`
		CourseXP       int      `json:"course_xp"`
		TopicIDs       []string `json:"topic_ids"`
	}

	DashboardData struct {
		Courses  []CourseSummaryData `json:"courses"`
		Progress []ProgressData      `json:"progress"`

		StreakDays           int                 `json:"streak_days"`
		LastPracticedSheetID string              `json:"last_practiced_sheet_id,omitempty"`
		ResumePractice       *ResumePracticeData `json:"resume_practice,omitempty"`
	}

	ResumePracticeData struct {
		SheetID    string `json:"sheet_id"`
		TopicID    string `json:"topic_id,omitempty"`
		TopicTitle string `json:"topic_title,omitempty"`
		Level      int    `json:"level"`
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

	xpByCourse := make(map[string]int)
	if xpEntries, xpErr := app.Repositories.StudentCourseXP.ListByStudent(ctx, studentID); xpErr != nil {
		log.Printf("[dashboard] could not load course xp student_id=%s err=%v", studentID, xpErr)
	} else {
		for _, entry := range xpEntries {
			xpByCourse[entry.CourseID] = entry.TotalXP
		}
	}

	resumePractice, err := app.Repositories.StudentPracticeState.GetLastOpenedPractice(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptGetError, err)
	}

	loc := studentLocation(ctx, app, studentID)

	courses := make([]CourseSummaryData, 0, len(summaries))
	for _, s := range summaries {
		courses = append(courses, CourseSummaryData{
			CourseID:       s.CourseID,
			SchoolID:       s.SchoolID,
			SchoolName:     s.SchoolName,
			Title:          s.Title,
			Subject:        s.Subject,
			GradeName:      s.GradeName,
			PracticeSheets: s.PracticeSheets,
			LevelTests:     s.LevelTests,
			Notebooks:      s.Notebooks,
			CurrentLevel:   s.CurrentLevel,
			CourseXP:       xpByCourse[s.CourseID],
			TopicIDs:       s.TopicIDs,
		})
	}

	progressData := make([]ProgressData, 0, len(progress))
	for _, p := range progress {
		progressData = append(progressData, toProgressData(p, loc))
	}
	streak := domain.CurrentStreak(progress, loc)

	data := DashboardData{
		Courses:    courses,
		Progress:   progressData,
		StreakDays: streak,
	}
	if resumePractice != nil {
		data.LastPracticedSheetID = resumePractice.SheetID
		data.ResumePractice = &ResumePracticeData{
			SheetID:    resumePractice.SheetID,
			TopicID:    resumePractice.TopicID,
			TopicTitle: resumePractice.TopicTitle,
			Level:      resumePractice.Level,
		}
	}
	return &DashboardOutput{Data: data}, nil
}

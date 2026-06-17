package studentreport

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GeneratePDFUsecase interface {
		Execute(ctx context.Context, teacherID string, isAdmin bool, filter domain.StudentReportFilter) ([]byte, apperrors.ApplicationError)
	}

	generatePDFUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewGeneratePDFUsecase(contextFactory appcontext.Factory) GeneratePDFUsecase {
	return &generatePDFUsecase{contextFactory: contextFactory}
}

func (u *generatePDFUsecase) Execute(ctx context.Context, teacherID string, isAdmin bool, filter domain.StudentReportFilter) ([]byte, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !isAdmin {
		hasAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, teacherID, filter.StudentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

	student, err := app.Repositories.UserProfile.Get(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("student not found")
	}

	var topicProgress []domain.StudentTopicProgress
	if filter.CourseID != "" {
		topicProgress, err = app.Repositories.StudentProgress.ListByStudentAndCourse(ctx, filter.StudentID, filter.CourseID)
	} else {
		topicProgress, err = app.Repositories.StudentProgress.ListByStudent(ctx, filter.StudentID)
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	if filter.From != nil || filter.To != nil {
		topicProgress = filterProgressByDate(topicProgress, filter.From, filter.To)
	}

	courseProgressList, err := app.Repositories.CourseProgress.ListByStudent(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	courseItems := buildCourseProgressItems(ctx, app, courseProgressList, topicProgress)
	if filter.CourseID != "" {
		filtered := make([]domain.CourseProgressItem, 0)
		for _, c := range courseItems {
			if c.CourseID == filter.CourseID {
				filtered = append(filtered, c)
			}
		}
		courseItems = filtered
	}

	summary := calculateSummary(topicProgress)

	recentAttempts := []domain.StudentAttempt{}
	dailyAttempts, err := app.Repositories.StudentAttempt.GetDailyAttempts(ctx, filter.StudentID, filter.CourseID, filter.From, filter.To)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	reportData := &domain.StudentReportData{
		Student:     *student,
		GeneratedAt: time.Now(),
		Period: domain.ReportPeriod{
			From: filter.From,
			To:   filter.To,
		},
		Summary:        summary,
		TopicProgress:  topicProgress,
		CourseProgress: courseItems,
		DailyAttempts:  dailyAttempts,
		RecentAttempts: recentAttempts,
	}

	builder := NewPDFBuilder(reportData)
	pdfBytes, err := builder.Build()
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return pdfBytes, nil
}

func filterProgressByDate(progress []domain.StudentTopicProgress, from, to *time.Time) []domain.StudentTopicProgress {
	filtered := make([]domain.StudentTopicProgress, 0)
	for _, p := range progress {
		if p.LastPracticedAt == nil {
			continue
		}
		if from != nil && p.LastPracticedAt.Before(*from) {
			continue
		}
		if to != nil && p.LastPracticedAt.After(*to) {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
}

func buildCourseProgressItems(ctx context.Context, app *appcontext.Context, courseProgress []domain.StudentCourseProgress, topicProgress []domain.StudentTopicProgress) []domain.CourseProgressItem {
	courseTopics := make(map[string][]domain.StudentTopicProgress)

	for _, tp := range topicProgress {
		topic, err := app.Repositories.Topic.Get(ctx, tp.TopicID)
		if err != nil || topic == nil {
			continue
		}
		courseTopics[topic.CourseID] = append(courseTopics[topic.CourseID], tp)
	}

	items := make([]domain.CourseProgressItem, 0, len(courseProgress))
	for _, cp := range courseProgress {
		course, err := app.Repositories.Course.Get(ctx, cp.CourseID)
		if err != nil || course == nil {
			continue
		}

		// Get topics for this course
		topics := courseTopics[cp.CourseID]
		topicCount := len(topics)
		avgMastery := 0.0
		var lastActivity *time.Time

		if topicCount > 0 {
			var total float64
			for _, t := range topics {
				total += t.MasteryScore
				if lastActivity == nil || (t.LastPracticedAt != nil && t.LastPracticedAt.After(*lastActivity)) {
					lastActivity = t.LastPracticedAt
				}
			}
			avgMastery = total / float64(topicCount)
		}

		items = append(items, domain.CourseProgressItem{
			CourseID:       cp.CourseID,
			CourseTitle:    course.Title,
			CurrentLevel:   cp.CurrentLevel,
			TopicCount:     topicCount,
			AverageMastery: avgMastery,
			LastActivity:   lastActivity,
		})
	}

	return items
}

func calculateSummary(progress []domain.StudentTopicProgress) domain.ReportSummary {
	if len(progress) == 0 {
		return domain.ReportSummary{}
	}

	var totalMastery float64
	var totalAttempts, correctAttempts, maxStreak int

	for _, p := range progress {
		totalMastery += p.MasteryScore
		totalAttempts += p.TotalAttempts
		correctAttempts += p.CorrectAttempts
		if p.StreakDays > maxStreak {
			maxStreak = p.StreakDays
		}
	}

	accuracyRate := 0.0
	if totalAttempts > 0 {
		accuracyRate = float64(correctAttempts) / float64(totalAttempts) * 100
	}

	return domain.ReportSummary{
		TopicsPracticed: len(progress),
		AverageMastery:  totalMastery / float64(len(progress)),
		TotalAttempts:   totalAttempts,
		CorrectAttempts: correctAttempts,
		AccuracyRate:    accuracyRate,
		CurrentStreak:   maxStreak,
	}
}

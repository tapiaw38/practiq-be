package studentreport

import (
	"context"
	"sort"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GeneratePDFUsecase interface {
	Execute(ctx context.Context, teacherID string, filter domain.StudentReportFilter) ([]byte, apperrors.ApplicationError)
}

type generatePDFUsecase struct {
	factory appcontext.Factory
}

func NewGeneratePDFUsecase(factory appcontext.Factory) GeneratePDFUsecase {
	return &generatePDFUsecase{factory: factory}
}

func (u *generatePDFUsecase) Execute(ctx context.Context, teacherID string, filter domain.StudentReportFilter) ([]byte, apperrors.ApplicationError) {
	app := u.factory()

	// Get student profile
	student, err := app.Repositories.UserProfile.Get(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("student not found")
	}

	// Get topic progress
	var topicProgress []domain.StudentTopicProgress
	if filter.CourseID != "" {
		topicProgress, err = app.Repositories.StudentProgress.ListByStudentAndCourse(ctx, filter.StudentID, filter.CourseID)
	} else {
		topicProgress, err = app.Repositories.StudentProgress.ListByStudent(ctx, filter.StudentID)
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	// Filter by date if provided
	if filter.From != nil || filter.To != nil {
		topicProgress = filterProgressByDate(topicProgress, filter.From, filter.To)
	}

	// Get course progress
	courseProgressList, err := app.Repositories.CourseProgress.ListByStudent(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	// Build course progress items with additional data
	courseItems := buildCourseProgressItems(ctx, app, courseProgressList, topicProgress)

	// Filter course items if specific course requested
	if filter.CourseID != "" {
		filtered := make([]domain.CourseProgressItem, 0)
		for _, c := range courseItems {
			if c.CourseID == filter.CourseID {
				filtered = append(filtered, c)
			}
		}
		courseItems = filtered
	}

	// Calculate summary
	summary := calculateSummary(topicProgress)

	// Get recent attempts (we'll need to aggregate from multiple sheets)
	// For simplicity, get attempts from last 30 days
	recentAttempts := []domain.StudentAttempt{}

	// Calculate daily attempts from topic progress (approximation)
	dailyAttempts := calculateDailyAttempts(topicProgress, filter.From, filter.To)

	// Build report data
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

	// Generate PDF
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
	// Group topic progress by course - we'll calculate from course progress data
	courseTopics := make(map[string][]domain.StudentTopicProgress)
	_ = courseTopics // Used below when we have topic-course mapping

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

func calculateDailyAttempts(progress []domain.StudentTopicProgress, from, to *time.Time) []domain.DailyAttemptCount {
	// This is an approximation since we don't have daily breakdown
	// In a real implementation, we'd query student_attempts grouped by day
	dailyMap := make(map[string]domain.DailyAttemptCount)

	for _, p := range progress {
		if p.LastPracticedAt == nil {
			continue
		}
		dateKey := p.LastPracticedAt.Format("2006-01-02")
		existing := dailyMap[dateKey]
		existing.Date = time.Date(p.LastPracticedAt.Year(), p.LastPracticedAt.Month(), p.LastPracticedAt.Day(), 0, 0, 0, 0, p.LastPracticedAt.Location())
		existing.Total += p.TotalAttempts
		existing.Correct += p.CorrectAttempts
		dailyMap[dateKey] = existing
	}

	result := make([]domain.DailyAttemptCount, 0, len(dailyMap))
	for _, v := range dailyMap {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	// Limit to last 30 days
	if len(result) > 30 {
		result = result[len(result)-30:]
	}

	return result
}

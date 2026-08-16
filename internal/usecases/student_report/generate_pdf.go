package studentreport

import (
	"context"
	"sort"
	"time"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
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
		if filter.CourseID != "" {
			course, err := app.Repositories.Course.Get(ctx, filter.CourseID)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
			}
			if course == nil || course.TeacherID != teacherID {
				return nil, apperrors.NewForbiddenError()
			}
		}
	}

	student, err := app.Repositories.UserProfile.Get(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("student not found")
	}

	// Sharing one course with the student is enough to pass HasAccess, but the
	// unfiltered report used to load every topic row — including courses that
	// belong to other teachers. Without a course filter, restrict it to the
	// requester's own courses.
	var ownCourseIDs map[string]bool
	if !isAdmin && filter.CourseID == "" {
		ownCourses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{TeacherID: teacherID})
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
		}
		ownCourseIDs = make(map[string]bool, len(ownCourses))
		for _, c := range ownCourses {
			ownCourseIDs[c.ID] = true
		}
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
	if ownCourseIDs != nil {
		topicProgress = filterProgressByCourses(ctx, app, topicProgress, ownCourseIDs)
	}

	if filter.From != nil || filter.To != nil {
		topicProgress = filterProgressByDate(topicProgress, filter.From, filter.To)
	}

	courseProgressList, err := app.Repositories.CourseProgress.ListByStudent(ctx, filter.StudentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}
	if ownCourseIDs != nil {
		kept := courseProgressList[:0]
		for _, cp := range courseProgressList {
			if ownCourseIDs[cp.CourseID] {
				kept = append(kept, cp)
			}
		}
		courseProgressList = kept
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

	recentAttempts := []domain.StudentAttempt{}
	var dailyAttempts []domain.DailyAttemptCount
	if ownCourseIDs != nil {
		dailyAttempts, err = dailyAttemptsForCourses(ctx, app, filter.StudentID, ownCourseIDs, filter.From, filter.To)
	} else {
		dailyAttempts, err = app.Repositories.StudentAttempt.GetDailyAttempts(ctx, filter.StudentID, filter.CourseID, filter.From, filter.To)
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	summary := calculateSummary(topicProgress, dailyAttempts)

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

func dailyAttemptsForCourses(
	ctx context.Context,
	app *appcontext.Context,
	studentID string,
	courseIDs map[string]bool,
	from, to *time.Time,
) ([]domain.DailyAttemptCount, error) {
	byDate := make(map[string]domain.DailyAttemptCount)
	for courseID := range courseIDs {
		attempts, err := app.Repositories.StudentAttempt.GetDailyAttempts(ctx, studentID, courseID, from, to)
		if err != nil {
			return nil, err
		}
		for _, attempt := range attempts {
			key := attempt.Date.Format("2006-01-02")
			current := byDate[key]
			current.Date = attempt.Date
			current.Total += attempt.Total
			current.Correct += attempt.Correct
			byDate[key] = current
		}
	}

	result := make([]domain.DailyAttemptCount, 0, len(byDate))
	for _, attempt := range byDate {
		result = append(result, attempt)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})
	return result, nil
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
	// Batch load all topics
	topicIDs := make([]string, 0, len(topicProgress))
	for _, tp := range topicProgress {
		topicIDs = append(topicIDs, tp.TopicID)
	}

	topics, err := app.Repositories.Topic.GetByIDs(ctx, topicIDs)
	topicMap := make(map[string]*domain.Topic)
	if err == nil {
		for i := range topics {
			topicMap[topics[i].ID] = &topics[i]
		}
	}

	// Build course -> topics mapping
	courseTopics := make(map[string][]domain.StudentTopicProgress)
	for _, tp := range topicProgress {
		if topic, ok := topicMap[tp.TopicID]; ok {
			courseTopics[topic.CourseID] = append(courseTopics[topic.CourseID], tp)
		}
	}

	// Batch load all courses
	courseIDs := make([]string, 0, len(courseProgress))
	for _, cp := range courseProgress {
		courseIDs = append(courseIDs, cp.CourseID)
	}

	courses, err := app.Repositories.Course.GetByIDs(ctx, courseIDs)
	courseMap := make(map[string]*domain.Course)
	if err == nil {
		for i := range courses {
			courseMap[courses[i].ID] = &courses[i]
		}
	}

	// Build items
	items := make([]domain.CourseProgressItem, 0, len(courseProgress))
	for _, cp := range courseProgress {
		course, ok := courseMap[cp.CourseID]
		if !ok {
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

func calculateSummary(progress []domain.StudentTopicProgress, dailyAttempts []domain.DailyAttemptCount) domain.ReportSummary {
	var totalMastery float64
	var maxStreak int

	for _, p := range progress {
		totalMastery += p.MasteryScore
		if p.StreakDays > maxStreak {
			maxStreak = p.StreakDays
		}
	}

	var totalAttempts, correctAttempts int
	for _, da := range dailyAttempts {
		totalAttempts += da.Total
		correctAttempts += da.Correct
	}

	accuracyRate := 0.0
	if totalAttempts > 0 {
		accuracyRate = float64(correctAttempts) / float64(totalAttempts) * 100
	}

	avgMastery := 0.0
	if len(progress) > 0 {
		avgMastery = totalMastery / float64(len(progress))
	}

	return domain.ReportSummary{
		TopicsPracticed: len(progress),
		AverageMastery:  avgMastery,
		TotalAttempts:   totalAttempts,
		CorrectAttempts: correctAttempts,
		AccuracyRate:    accuracyRate,
		CurrentStreak:   maxStreak,
	}
}

// filterProgressByCourses keeps only the topics that belong to the given
// courses. Topic rows carry no course id, so the topic is resolved to find it.
func filterProgressByCourses(
	ctx context.Context,
	app *appcontext.Context,
	progress []domain.StudentTopicProgress,
	allowed map[string]bool,
) []domain.StudentTopicProgress {
	kept := make([]domain.StudentTopicProgress, 0, len(progress))
	for _, p := range progress {
		topic, err := app.Repositories.Topic.Get(ctx, p.TopicID)
		// A topic that cannot be resolved (deleted course) is dropped: the
		// safe default for a report that crosses teachers is to omit.
		if err != nil || topic == nil || !allowed[topic.CourseID] {
			continue
		}
		kept = append(kept, p)
	}
	return kept
}

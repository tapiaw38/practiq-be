package practicesheet

import (
	"context"
	"fmt"
	"log"
	"time"

	enrollmentRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

const sheetTypeLevelTest = "level_test"

func isCourseTeacher(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) bool {
	_, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, courseID)
	return appErr == nil
}

func ensureSheetIsOpen(ctx context.Context, app *appcontext.Context, ps *domain.PracticeSheet, requesterID string, isSuperAdmin bool) apperrors.ApplicationError {
	if isCourseTeacher(ctx, app, requesterID, isSuperAdmin, ps.CourseID) {
		return nil
	}
	if ps.SheetType == sheetTypeLevelTest {
		progress, err := app.Repositories.CourseProgress.Get(ctx, requesterID, ps.CourseID)
		if err != nil {
			return apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
		}
		currentLevel := 1
		if progress != nil {
			currentLevel = progress.CurrentLevel
		}
		if ps.Level != currentLevel {
			return apperrors.NewBadRequestError("this level test does not belong to the student's current level")
		}
	}
	if ps.ScheduledAt == nil {
		return nil
	}

	switch sheetWindowState(ps, time.Now()) {
	case windowNotYetOpen:
		return apperrors.NewApplicationError(mappings.PracticeSheetNotYetAvailableError, nil)
	case windowClosed:
		return apperrors.NewApplicationError(mappings.PracticeSheetExpiredError, nil)
	default:
		return nil
	}
}

type windowState int

const (
	windowOpen windowState = iota
	windowNotYetOpen
	windowClosed
)

func sheetWindowState(ps *domain.PracticeSheet, now time.Time) windowState {
	if ps.ScheduledAt == nil {
		return windowOpen
	}
	if now.Before(*ps.ScheduledAt) {
		return windowNotYetOpen
	}
	if ps.AvailableUntil != nil && !now.Before(*ps.AvailableUntil) {
		return windowClosed
	}
	return windowOpen
}

func notifyScheduledLevelTest(ctx context.Context, app *appcontext.Context, ps domain.PracticeSheet) {
	if ps.SheetType != sheetTypeLevelTest || ps.ScheduledAt == nil {
		if err := app.Repositories.Notification.DeleteByResource(ctx, domain.NotificationLevelTestScheduled, ps.ID); err != nil {
			log.Printf("[practice_sheet] failed to clear notifications sheet_id=%s err=%v", ps.ID, err)
		}
		return
	}

	students, err := app.Repositories.Enrollment.ListStudents(ctx, enrollmentRepo.ListFilter{CourseID: ps.CourseID})
	if err != nil {
		log.Printf("[practice_sheet] failed to list students for notifications course_id=%s err=%v", ps.CourseID, err)
		return
	}

	courseTitle := ""
	if course, err := app.Repositories.Course.Get(ctx, ps.CourseID); err == nil && course != nil {
		courseTitle = course.Title
	}

	body := "Tenés una prueba de nivel programada"
	if courseTitle != "" {
		body = fmt.Sprintf("%s en %s", body, courseTitle)
	}

	for _, student := range students {
		notification := domain.Notification{
			UserID:       student.ID,
			Type:         domain.NotificationLevelTestScheduled,
			Title:        ps.Title,
			Body:         body,
			ResourceType: domain.NotificationResourcePracticeSheet,
			ResourceID:   ps.ID,
			ScheduledAt:  ps.ScheduledAt,
		}
		if err := app.Repositories.Notification.Upsert(ctx, notification); err != nil {
			log.Printf("[practice_sheet] failed to notify student_id=%s sheet_id=%s err=%v", student.ID, ps.ID, err)
		}
	}
}

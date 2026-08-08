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
)

const sheetTypeLevelTest = "level_test"

// isCourseTeacher reports whether the requester owns the course. Teachers and
// admins bypass the schedule so they can review a test before its date.
func isCourseTeacher(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, courseID string) bool {
	if isAdmin {
		return true
	}
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil || course == nil {
		return false
	}
	return course.TeacherID == requesterID
}

// ensureSheetIsOpen blocks students before and after a scheduled level test.
// Teachers and admins bypass the schedule so they can review or reprogram it.
func ensureSheetIsOpen(ctx context.Context, app *appcontext.Context, ps *domain.PracticeSheet, requesterID string, isAdmin bool) apperrors.ApplicationError {
	if isCourseTeacher(ctx, app, requesterID, isAdmin, ps.CourseID) {
		return nil
	}
	if ps.ScheduledAt == nil {
		return nil
	}
	if time.Now().Before(*ps.ScheduledAt) {
		return apperrors.NewApplicationError(mappings.PracticeSheetNotYetAvailableError, nil)
	}
	return apperrors.NewApplicationError(mappings.PracticeSheetExpiredError, nil)
}

// notifyScheduledLevelTest tells every enrolled student when a level test gets a
// date. Unscheduling or switching the sheet back to practice removes the
// notification instead.
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

	// No date in the body on purpose: scheduled_at travels structured so the
	// client renders it in the student's own timezone.
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

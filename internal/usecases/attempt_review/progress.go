package attemptreview

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

const (
	sheetTypeLevelTest     = "level_test"
	levelTestPassThreshold = 75.0
)

func applyLevelTestOutcome(ctx context.Context, app *appcontext.Context, attemptID string) {
	attemptCtx, err := app.Repositories.StudentAttempt.GetAttemptContext(ctx, attemptID)
	if err != nil {
		log.Printf("[attempt_review] could not read attempt context attempt_id=%s err=%v", attemptID, err)
		return
	}
	if attemptCtx.SheetType != sheetTypeLevelTest || attemptCtx.PracticeSheetID == "" {
		return
	}

	outcome, err := app.Repositories.StudentAttempt.GetSheetOutcome(ctx, attemptCtx.StudentID, attemptCtx.PracticeSheetID)
	if err != nil {
		log.Printf("[attempt_review] could not read sheet outcome sheet_id=%s err=%v", attemptCtx.PracticeSheetID, err)
		return
	}

	if outcome.Pending > 0 || outcome.Total == 0 {
		return
	}

	score := float64(outcome.Correct) / float64(outcome.Total) * 100
	if score < levelTestPassThreshold {
		return
	}

	topicID := attemptCtx.SheetTopicID
	if topicID == "" {
		topicID = attemptCtx.ExerciseTopicID
	}
	if topicID == "" {
		return
	}

	sheet, err := app.Repositories.PracticeSheet.Get(ctx, attemptCtx.PracticeSheetID)
	if err != nil || sheet == nil {
		log.Printf("[attempt_review] could not read sheet sheet_id=%s err=%v", attemptCtx.PracticeSheetID, err)
		return
	}

	progress, err := app.Repositories.StudentProgress.Get(ctx, attemptCtx.StudentID, topicID)
	if err != nil {
		log.Printf("[attempt_review] could not read progress student_id=%s err=%v", attemptCtx.StudentID, err)
		return
	}

	currentLevel := 1
	if progress != nil {
		currentLevel = progress.CurrentLevel
	}

	targetLevel := sheet.Level + 1

	advanceCourseProgress(ctx, app, attemptCtx.StudentID, sheet.CourseID, targetLevel)
	if currentLevel >= targetLevel {
		log.Printf("[attempt_review] promotion already applied student_id=%s sheet_id=%s level=%d",
			attemptCtx.StudentID, attemptCtx.PracticeSheetID, currentLevel)
		return
	}

	updated := domain.StudentTopicProgress{
		StudentID:    attemptCtx.StudentID,
		TopicID:      topicID,
		MasteryScore: score,
		CurrentLevel: targetLevel,
	}
	if progress != nil {
		updated.TotalAttempts = progress.TotalAttempts
		updated.CorrectAttempts = progress.CorrectAttempts
		updated.StreakDays = progress.StreakDays
	}

	if err := app.Repositories.StudentProgress.Upsert(ctx, updated); err != nil {
		log.Printf("[attempt_review] could not save progress student_id=%s err=%v", attemptCtx.StudentID, err)
		return
	}
	log.Printf("[attempt_review] level test passed after review student_id=%s sheet_id=%s score=%.0f level=%d",
		attemptCtx.StudentID, attemptCtx.PracticeSheetID, score, updated.CurrentLevel)
}

func advanceCourseProgress(ctx context.Context, app *appcontext.Context, studentID, courseID string, targetLevel int) {
	if courseID == "" {
		return
	}
	courseProgress, err := app.Repositories.CourseProgress.Get(ctx, studentID, courseID)
	if err != nil {
		log.Printf("[attempt_review] could not read course progress student_id=%s err=%v", studentID, err)
		return
	}
	if courseProgress != nil && courseProgress.CurrentLevel >= targetLevel {
		return
	}
	if err := app.Repositories.CourseProgress.Upsert(ctx, studentID, courseID, targetLevel); err != nil {
		log.Printf("[attempt_review] could not save course progress student_id=%s err=%v", studentID, err)
	}
}

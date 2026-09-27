package practicesheet

import (
	"context"
	"log"
	"strings"

	studentcoursexp "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_course_xp"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

const (
	xpAttempt          = 1
	xpCorrect          = 10
	xpPracticeComplete = 10
	xpLevelTestPass    = 40
)

type xpExercise struct {
	ExerciseID string
	Correct    bool
	Ungraded   bool
}

type XPBreakdownEntry struct {
	EventType string `json:"event_type"`
	Points    int    `json:"points"`
	Count     int    `json:"count"`
}

type xpOutcome struct {
	Gained    int
	Balance   int
	Breakdown []XPBreakdownEntry
}

func awardPracticeXP(ctx context.Context, app *appcontext.Context, studentID string, sheet *domain.PracticeSheet, attempts []AttemptInput, exercises []xpExercise, levelTestPassed bool) (result xpOutcome) {

	balanceKnown := false
	defer func() {
		if balanceKnown {
			return
		}
		current, err := app.Repositories.StudentCourseXP.Balance(ctx, studentID, sheet.CourseID)
		if err != nil {
			log.Printf("[practice_xp] balance read failed student_id=%s course_id=%s err=%v", studentID, sheet.CourseID, err)
			return
		}
		result.Balance = current
	}()

	breakdownIndex := map[string]int{}
	countIn := func(eventType string, points int) {
		if position, ok := breakdownIndex[eventType]; ok {
			result.Breakdown[position].Points += points
			result.Breakdown[position].Count++
			return
		}
		breakdownIndex[eventType] = len(result.Breakdown)
		result.Breakdown = append(result.Breakdown, XPBreakdownEntry{EventType: eventType, Points: points, Count: 1})
	}

	worked := attemptedExerciseIDs(attempts)
	award := func(eventKey, eventType string, points int, exerciseID string) {
		outcome, err := app.Repositories.StudentCourseXP.Award(ctx, studentcoursexp.AwardInput{
			StudentID: studentID, CourseID: sheet.CourseID, EventKey: eventKey, EventType: eventType,
			Points: points, PracticeSheetID: sheet.ID, ExerciseID: exerciseID,
		})
		if err != nil {

			log.Printf("[practice_xp] award failed student_id=%s course_id=%s event=%s err=%v", studentID, sheet.CourseID, eventKey, err)
			return
		}
		if outcome.Awarded {
			result.Gained += points
			countIn(eventType, points)
		}
		result.Balance = outcome.TotalXP
		balanceKnown = true
	}

	for _, exercise := range exercises {
		if !worked[exercise.ExerciseID] {
			continue
		}
		award("exercise_attempt:"+exercise.ExerciseID, "exercise_attempt", xpAttempt, exercise.ExerciseID)
		if exercise.Correct && !exercise.Ungraded {
			award("exercise_correct:"+exercise.ExerciseID, "exercise_correct", xpCorrect, exercise.ExerciseID)
		}
	}

	if sheet.SheetType == sheetTypeLevelTest {
		if levelTestPassed {
			award("level_test_pass:"+sheet.ID, "level_test_pass", xpLevelTestPass, "")
		}
		return result
	}
	if allSheetExercisesAttempted(sheet.Exercises, worked) {
		award("practice_complete:"+sheet.ID, "practice_complete", xpPracticeComplete, "")
	}
	return result
}

func attemptedExerciseIDs(attempts []AttemptInput) map[string]bool {
	result := make(map[string]bool, len(attempts))
	for _, attempt := range attempts {
		if strings.TrimSpace(attempt.AnswerText) != "" || strings.TrimSpace(attempt.CanvasData) != "" || strings.TrimSpace(attempt.AttachmentURL) != "" {
			result[attempt.ExerciseID] = true
		}
	}
	return result
}

func allSheetExercisesAttempted(exercises []domain.PracticeSheetExercise, attempted map[string]bool) bool {
	if len(exercises) == 0 {
		return false
	}
	for _, exercise := range exercises {
		if !attempted[exercise.Exercise.ID] {
			return false
		}
	}
	return true
}

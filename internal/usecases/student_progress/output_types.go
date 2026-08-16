package studentprogress

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type (
	ProgressData struct {
		TopicID         string  `json:"topic_id"`
		TopicTitle      string  `json:"topic_title"`
		StrategyID      string  `json:"strategy_id"`
		MasteryScore    float64 `json:"mastery_score"`
		CurrentLevel    int     `json:"current_level"`
		TotalAttempts   int     `json:"total_attempts"`
		CorrectAttempts int     `json:"correct_attempts"`
		StreakDays      int     `json:"streak_days"`
		LastPracticedAt string  `json:"last_practiced_at"`
	}

	AttemptData struct {
		ID              string  `json:"id"`
		StudentID       string  `json:"student_id"`
		ExerciseID      string  `json:"exercise_id"`
		PracticeSheetID string  `json:"practice_sheet_id"`
		AnswerText      string  `json:"answer_text"`
		ImageURL        string  `json:"image_url,omitempty"`
		AIFeedback      string  `json:"ai_feedback,omitempty"`
		IsCorrect       bool    `json:"is_correct"`
		Score           float64 `json:"score"`
		TimeSpentSecs   int     `json:"time_spent_seconds"`
		HintsUsed       int     `json:"hints_used"`
		CreatedAt       string  `json:"created_at"`
	}
)

func toProgressData(p domain.StudentTopicProgress, loc *time.Location) ProgressData {
	lastPracticed := ""
	if p.LastPracticedAt != nil {
		lastPracticed = p.LastPracticedAt.Format("2006-01-02T15:04:05Z")
	}
	return ProgressData{
		TopicID:         p.TopicID,
		TopicTitle:      p.TopicTitle,
		StrategyID:      p.StrategyID,
		MasteryScore:    p.MasteryScore,
		CurrentLevel:    p.CurrentLevel,
		TotalAttempts:   p.TotalAttempts,
		CorrectAttempts: p.CorrectAttempts,
		StreakDays:      domain.EffectiveStreak(p, loc),
		LastPracticedAt: lastPracticed,
	}
}

func toAttemptData(a domain.StudentAttempt) AttemptData {
	return AttemptData{
		ID:              a.ID,
		StudentID:       a.StudentID,
		ExerciseID:      a.ExerciseID,
		PracticeSheetID: a.PracticeSheetID,
		AnswerText:      a.AnswerText,
		ImageURL:        a.ImageURL,
		AIFeedback:      a.AIFeedback,
		IsCorrect:       a.IsCorrect,
		Score:           a.Score,
		TimeSpentSecs:   a.TimeSpentSecs,
		HintsUsed:       a.HintsUsed,
		CreatedAt:       a.CreatedAt.Format(time.RFC3339),
	}
}

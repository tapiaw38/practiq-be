package studentattempt

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, a domain.StudentAttempt) (string, error) {
	query := `
		INSERT INTO student_attempts (student_id, exercise_id, practice_sheet_id, answer_text, image_url, ai_feedback, is_correct, score, time_spent_seconds, hints_used)
		VALUES ($1, $2, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, a.StudentID, a.ExerciseID, a.PracticeSheetID, a.AnswerText, a.ImageURL, a.AIFeedback, a.IsCorrect, a.Score, a.TimeSpentSecs, a.HintsUsed).Scan(&id)
	return id, err
}

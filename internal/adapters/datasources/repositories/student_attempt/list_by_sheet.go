package studentattempt

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListBySheet(ctx context.Context, studentID, sheetID string) ([]domain.StudentAttempt, error) {
	query := `
		SELECT id, student_id, exercise_id, COALESCE(practice_sheet_id::text,''), answer_text, COALESCE(image_url,''), COALESCE(ai_feedback,''), is_correct, score, time_spent_seconds, hints_used, created_at
		FROM student_attempts
		WHERE student_id = $1 AND practice_sheet_id = $2::uuid
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID, sheetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []domain.StudentAttempt
	for rows.Next() {
		var a domain.StudentAttempt
		if err := rows.Scan(&a.ID, &a.StudentID, &a.ExerciseID, &a.PracticeSheetID, &a.AnswerText, &a.ImageURL, &a.AIFeedback, &a.IsCorrect, &a.Score, &a.TimeSpentSecs, &a.HintsUsed, &a.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

package studentattempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.StudentAttempt) (string, error)
	ListBySheet(ctx context.Context, studentID, sheetID string) ([]domain.StudentAttempt, error)
	SaveCanvasWork(ctx context.Context, attemptID, imageData string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, a domain.StudentAttempt) (string, error) {
	query := `
		INSERT INTO student_attempts (student_id, exercise_id, practice_sheet_id, answer_text, ai_feedback, is_correct, score, time_spent_seconds, hints_used)
		VALUES ($1, $2, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, a.StudentID, a.ExerciseID, a.PracticeSheetID, a.AnswerText, a.AIFeedback, a.IsCorrect, a.Score, a.TimeSpentSecs, a.HintsUsed).Scan(&id)
	return id, err
}

func (r *repository) SaveCanvasWork(ctx context.Context, attemptID, imageData string) error {
	updateQuery := `
		UPDATE student_work_canvas
		SET image_url = $2
		WHERE attempt_id = $1::uuid
	`
	result, err := r.db.ExecContext(ctx, updateQuery, attemptID, imageData)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO student_work_canvas (attempt_id, image_url)
		VALUES ($1::uuid, $2)
	`
	_, err = r.db.ExecContext(ctx, insertQuery, attemptID, imageData)
	return err
}

func (r *repository) ListBySheet(ctx context.Context, studentID, sheetID string) ([]domain.StudentAttempt, error) {
	query := `
		SELECT id, student_id, exercise_id, COALESCE(practice_sheet_id::text,''), answer_text, COALESCE(ai_feedback,''), is_correct, score, time_spent_seconds, hints_used, created_at
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
		if err := rows.Scan(&a.ID, &a.StudentID, &a.ExerciseID, &a.PracticeSheetID, &a.AnswerText, &a.AIFeedback, &a.IsCorrect, &a.Score, &a.TimeSpentSecs, &a.HintsUsed, &a.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

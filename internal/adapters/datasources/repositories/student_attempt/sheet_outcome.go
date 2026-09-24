package studentattempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// GetSheetOutcome summarises a student's latest answer per exercise for a
// sheet. Attempts accumulate across submissions, so only the newest one per
// exercise counts — that is the submission the teacher is correcting.
func (r *repository) GetSheetOutcome(ctx context.Context, studentID, sheetID string) (domain.SheetOutcome, error) {
	var outcome domain.SheetOutcome
	err := r.db.QueryRowContext(ctx, `
		WITH latest AS (
			SELECT DISTINCT ON (exercise_id) exercise_id, is_correct, needs_teacher_review, teacher_reviewed_at
			FROM student_attempts
			WHERE student_id = $1 AND practice_sheet_id = $2::uuid
			ORDER BY exercise_id, created_at DESC
		)
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE is_correct),
			COUNT(*) FILTER (WHERE needs_teacher_review AND teacher_reviewed_at IS NULL)
		FROM latest
	`, studentID, sheetID).Scan(&outcome.Total, &outcome.Correct, &outcome.Pending)
	if err == sql.ErrNoRows {
		return outcome, nil
	}
	return outcome, err
}

// GetAttemptContext resolves the sheet and student an attempt belongs to, so a
// review can recompute that sheet's outcome.
func (r *repository) GetAttemptContext(ctx context.Context, attemptID string) (domain.AttemptContext, error) {
	var attemptCtx domain.AttemptContext
	err := r.db.QueryRowContext(ctx, `
		SELECT sa.student_id,
		       COALESCE(sa.practice_sheet_id::text, ''),
		       COALESCE(ps.sheet_type, ''),
		       COALESCE(ps.topic_id::text, ''),
		       COALESCE(e.topic_id::text, '')
		FROM student_attempts sa
		LEFT JOIN practice_sheets ps ON ps.id = sa.practice_sheet_id
		JOIN exercises e ON e.id = sa.exercise_id
		WHERE sa.id = $1
	`, attemptID).Scan(
		&attemptCtx.StudentID,
		&attemptCtx.PracticeSheetID,
		&attemptCtx.SheetType,
		&attemptCtx.SheetTopicID,
		&attemptCtx.ExerciseTopicID,
	)
	if err == sql.ErrNoRows {
		return attemptCtx, nil
	}
	return attemptCtx, err
}

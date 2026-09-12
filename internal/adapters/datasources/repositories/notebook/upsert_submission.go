package notebook

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// ErrStaleSubmission reports a delivery the student has already superseded.
//
// It is not a failure the caller should retry or surface: the newer answer is
// already stored, and this one arriving late changes nothing.
var ErrStaleSubmission = fmt.Errorf("a newer submission for this page is already stored")

func (r *repository) UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO notebook_submissions (page_id, student_id, canvas_data, answer_text, ai_recognized_text, ai_is_correct, ai_feedback, ai_reviewed_at, needs_teacher_review, version)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		WHERE EXISTS (
			SELECT 1 FROM notebook_pages np
			JOIN notebooks n ON n.id = np.notebook_id
			WHERE np.id = $1 AND n.deleted_at IS NULL
		)
		ON CONFLICT (page_id, student_id)
		DO UPDATE SET
			canvas_data = EXCLUDED.canvas_data,
			answer_text = EXCLUDED.answer_text,
			ai_recognized_text = EXCLUDED.ai_recognized_text,
			ai_is_correct = EXCLUDED.ai_is_correct,
			ai_feedback = EXCLUDED.ai_feedback,
			ai_reviewed_at = EXCLUDED.ai_reviewed_at,
			needs_teacher_review = EXCLUDED.needs_teacher_review,
			version = EXCLUDED.version,
			-- Saving an unchanged page is idempotent and must not erase a
			-- teacher's review. A changed answer starts a new review cycle.
			teacher_is_correct = CASE WHEN notebook_submissions.canvas_data IS DISTINCT FROM EXCLUDED.canvas_data OR notebook_submissions.answer_text IS DISTINCT FROM EXCLUDED.answer_text THEN NULL ELSE notebook_submissions.teacher_is_correct END,
			teacher_feedback = CASE WHEN notebook_submissions.canvas_data IS DISTINCT FROM EXCLUDED.canvas_data OR notebook_submissions.answer_text IS DISTINCT FROM EXCLUDED.answer_text THEN '' ELSE notebook_submissions.teacher_feedback END,
			teacher_reviewed_at = CASE WHEN notebook_submissions.canvas_data IS DISTINCT FROM EXCLUDED.canvas_data OR notebook_submissions.answer_text IS DISTINCT FROM EXCLUDED.answer_text THEN NULL ELSE notebook_submissions.teacher_reviewed_at END,
			updated_at = NOW()
		-- Saving happens after the assistant replies, so a slow delivery can
		-- finish after a later one. Without this the old answer would land last
		-- and overwrite the newer one, and the student would watch their work
		-- revert on its own.
		WHERE EXCLUDED.version >= notebook_submissions.version
	`, s.PageID, s.StudentID, s.CanvasData, s.AnswerText, s.AIRecognizedText, s.AIIsCorrect, s.AIFeedback, s.AIReviewedAt, s.NeedsTeacherReview, s.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Two different reasons land here, and they are told apart because one
		// is an error and the other is the guard working. A row already holding
		// a newer version means this delivery was superseded; anything else
		// means the page is gone.
		var newerExists bool
		if err := r.db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM notebook_submissions
				WHERE page_id = $1 AND student_id = $2 AND version > $3
			)`, s.PageID, s.StudentID, s.Version).Scan(&newerExists); err != nil {
			return err
		}
		if newerExists {
			return ErrStaleSubmission
		}
		return fmt.Errorf("notebook page not found or notebook has been deleted")
	}

	return nil
}

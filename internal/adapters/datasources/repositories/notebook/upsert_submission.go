package notebook

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO notebook_submissions (page_id, student_id, canvas_data, answer_text, ai_recognized_text, ai_is_correct, ai_feedback, ai_reviewed_at, needs_teacher_review)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9
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
			teacher_is_correct = NULL,
			teacher_feedback = '',
			teacher_reviewed_at = NULL,
			updated_at = NOW()
	`, s.PageID, s.StudentID, s.CanvasData, s.AnswerText, s.AIRecognizedText, s.AIIsCorrect, s.AIFeedback, s.AIReviewedAt, s.NeedsTeacherReview)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notebook page not found or notebook has been deleted")
	}

	return nil
}

package notebook

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) UpdateSubmissionAIReview(ctx context.Context, id string, recognizedText string, isCorrect *bool, feedback string, needsTeacherReview bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET ai_recognized_text = $1, ai_is_correct = $2, ai_feedback = $3, ai_reviewed_at = NOW(), needs_teacher_review = $4, updated_at = NOW()
		WHERE id = $5 AND ($6 = '' OR school_id = NULLIF($6,'')::uuid)
	`, recognizedText, isCorrect, feedback, needsTeacherReview, id, tenant.SchoolID(ctx))
	return err
}

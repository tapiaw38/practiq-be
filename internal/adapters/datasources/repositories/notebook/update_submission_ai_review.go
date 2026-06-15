package notebook

import "context"

func (r *repository) UpdateSubmissionAIReview(ctx context.Context, id string, recognizedText string, isCorrect *bool, feedback string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET ai_recognized_text = $1, ai_is_correct = $2, ai_feedback = $3, ai_reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`, recognizedText, isCorrect, feedback, id)
	return err
}

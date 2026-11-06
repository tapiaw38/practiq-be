package notebook

import "context"

func (r *repository) UpdateSubmissionTeacherReview(ctx context.Context, id string, isCorrect bool, feedback string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET teacher_is_correct = $1, teacher_feedback = $2, teacher_reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`, isCorrect, feedback, id)
	return err
}

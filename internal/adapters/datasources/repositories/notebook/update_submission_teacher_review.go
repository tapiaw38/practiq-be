package notebook

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) UpdateSubmissionTeacherReview(ctx context.Context, id string, isCorrect bool, feedback string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET teacher_is_correct = $1, teacher_feedback = $2, teacher_reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3 AND ($4 = '' OR school_id = NULLIF($4,'')::uuid)
	`, isCorrect, feedback, id, tenant.SchoolID(ctx))
	return err
}

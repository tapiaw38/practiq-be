package learningstrategy

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) UnassignFromCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_learning_strategies WHERE id = $1 AND ($2 = '' OR school_id = NULLIF($2,'')::uuid)`, id, tenant.SchoolID(ctx))
	return err
}

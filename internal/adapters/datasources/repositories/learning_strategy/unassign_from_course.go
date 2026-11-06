package learningstrategy

import "context"

func (r *repository) UnassignFromCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_learning_strategies WHERE id = $1`, id)
	return err
}

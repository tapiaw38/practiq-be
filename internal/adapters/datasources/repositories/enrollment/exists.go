package enrollment

import "context"

func (r *repository) Exists(ctx context.Context, courseID, studentID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM enrollments WHERE course_id=$1 AND student_id=$2`, courseID, studentID).Scan(&count)
	return count > 0, err
}

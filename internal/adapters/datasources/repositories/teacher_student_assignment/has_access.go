package teacherstudentassignment

import "context"

func (r *repository) HasAccess(ctx context.Context, teacherID, studentID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM teacher_student_assignments
			WHERE teacher_id = $1 AND student_id = $2
			UNION
			SELECT 1 FROM enrollments e
			JOIN courses c ON c.id = e.course_id
			WHERE c.teacher_id = $1 AND e.student_id = $2 AND c.deleted_at IS NULL
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, studentID).Scan(&exists)
	return exists, err
}

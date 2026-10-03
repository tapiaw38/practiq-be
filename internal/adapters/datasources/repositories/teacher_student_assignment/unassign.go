package teacherstudentassignment

import "context"

func (r *repository) Unassign(ctx context.Context, teacherID, studentID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM teacher_student_assignments
		WHERE teacher_id = $1 AND student_id = $2
	`, teacherID, studentID)
	return err
}

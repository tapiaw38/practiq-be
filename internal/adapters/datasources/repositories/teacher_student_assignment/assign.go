package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Assign(ctx context.Context, assignment domain.TeacherStudentAssignment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO teacher_student_assignments (teacher_id, student_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (teacher_id, student_id) DO UPDATE SET status = EXCLUDED.status
	`, assignment.TeacherID, assignment.StudentID, assignment.Status)
	return err
}

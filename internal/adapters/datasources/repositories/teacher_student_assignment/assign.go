package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Assign(ctx context.Context, assignment domain.TeacherStudentAssignment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO teacher_student_assignments (teacher_id, student_id, status, school_id)
		VALUES ($1, $2, $3, NULLIF($4,'')::uuid)
		ON CONFLICT (teacher_id, student_id) DO UPDATE SET status = EXCLUDED.status, school_id = EXCLUDED.school_id
	`, assignment.TeacherID, assignment.StudentID, assignment.Status, tenant.SchoolID(ctx))
	return err
}

package teacherstudentassignment

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Unassign(ctx context.Context, teacherID, studentID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM teacher_student_assignments
		WHERE teacher_id = $1 AND student_id = $2 AND ($3 = '' OR school_id = NULLIF($3,'')::uuid)
	`, teacherID, studentID, tenant.SchoolID(ctx))
	return err
}

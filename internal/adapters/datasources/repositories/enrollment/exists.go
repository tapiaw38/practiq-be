package enrollment

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Exists(ctx context.Context, courseID, studentID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM enrollments WHERE course_id=$1 AND student_id=$2 AND ($3 = '' OR school_id = NULLIF($3, '')::uuid)`, courseID, studentID, tenant.SchoolID(ctx)).Scan(&count)
	return count > 0, err
}

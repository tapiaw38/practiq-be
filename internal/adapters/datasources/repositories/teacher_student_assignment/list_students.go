package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) ListStudents(ctx context.Context, filter ListFilter) ([]domain.UserProfile, error) {
	return r.listUsers(ctx, `
		SELECT up.id, up.profile_type, up.academic_status, up.created_at
		FROM user_profiles up
		JOIN teacher_student_assignments tsa ON tsa.student_id = up.id
		WHERE tsa.teacher_id = $1 AND tsa.status = 'active' AND ($2 = '' OR tsa.school_id = NULLIF($2,'')::uuid)
		ORDER BY up.id ASC
	`, filter.UserID, tenant.SchoolID(ctx), filter.Limit, filter.Offset)
}

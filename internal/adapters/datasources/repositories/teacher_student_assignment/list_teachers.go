package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListTeachers(ctx context.Context, filter ListFilter) ([]domain.UserProfile, error) {
	return r.listUsers(ctx, `
		SELECT up.id, up.name, up.email, up.profile_type, up.academic_status, up.assistant_base_url, up.assistant_api_key, up.created_at
		FROM user_profiles up
		JOIN teacher_student_assignments tsa ON tsa.teacher_id = up.id
		WHERE tsa.student_id = $1 AND tsa.status = 'active'
		ORDER BY up.name ASC
	`, filter.UserID, filter.Limit, filter.Offset)
}

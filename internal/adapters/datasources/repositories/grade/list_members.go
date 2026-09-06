package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) ListMembers(ctx context.Context, gradeID string) ([]domain.UserProfile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT up.id, up.profile_type, up.created_at
		FROM user_profiles up
		JOIN grade_memberships gm ON gm.user_id = up.id
		WHERE gm.grade_id = $1 AND ($2 = '' OR g.school_id = NULLIF($2, '')::uuid)
		ORDER BY up.id ASC
	`, gradeID, tenant.SchoolID(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserProfile{}
	for rows.Next() {
		var user domain.UserProfile
		if err := rows.Scan(
			&user.ID,
			&user.ProfileType,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListMembers(ctx context.Context, gradeID string) ([]domain.UserProfile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT up.id, up.name, up.email, up.profile_type, up.created_at
		FROM user_profiles up
		JOIN grade_memberships gm ON gm.user_id = up.id
		WHERE gm.grade_id = $1
		ORDER BY up.name ASC
	`, gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserProfile{}
	for rows.Next() {
		var user domain.UserProfile
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProfileType,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

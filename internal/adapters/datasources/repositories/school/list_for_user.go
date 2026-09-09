package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListForUser(ctx context.Context, userID string) ([]domain.SchoolMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT school_id, user_id, role, active
		FROM school_members
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []domain.SchoolMember{}
	for rows.Next() {
		var m domain.SchoolMember
		if err := rows.Scan(&m.SchoolID, &m.UserID, &m.Role, &m.Active); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

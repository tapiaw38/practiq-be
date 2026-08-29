package grade

import (
	"context"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListGradesByUsers(ctx context.Context, userIDs []string) (map[string][]domain.Grade, error) {
	result := make(map[string][]domain.Grade, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT gm.user_id, g.id, g.name, COALESCE(g.description, ''), g.visual_theme, g.created_by, g.created_at
		FROM grades g
		JOIN grade_memberships gm ON gm.grade_id = g.id
		WHERE gm.user_id = ANY($1)
		ORDER BY gm.user_id, g.name ASC
	`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		var grade domain.Grade
		if err := rows.Scan(&userID, &grade.ID, &grade.Name, &grade.Description, &grade.VisualTheme, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		result[userID] = append(result[userID], grade)
	}

	return result, rows.Err()
}

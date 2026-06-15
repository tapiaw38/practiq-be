package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListUserGrades(ctx context.Context, userID string) ([]domain.Grade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id, g.name, COALESCE(g.description, ''), g.created_by, g.created_at
		FROM grades g
		JOIN grade_memberships gm ON gm.grade_id = g.id
		WHERE gm.user_id = $1
		ORDER BY g.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []domain.Grade{}
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context) ([]domain.Grade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description, ''), visual_theme, created_by, created_at
		FROM grades
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []domain.Grade{}
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.VisualTheme, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

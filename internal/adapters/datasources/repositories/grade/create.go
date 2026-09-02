package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, grade domain.Grade) (string, error) {
	query := `
		INSERT INTO grades (name, description, visual_theme, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, grade.Name, grade.Description, grade.VisualTheme, grade.CreatedBy).Scan(&id)
	return id, err
}

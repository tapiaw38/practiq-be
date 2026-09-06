package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Create(ctx context.Context, grade domain.Grade) (string, error) {
	query := `
		INSERT INTO grades (name, description, visual_theme, created_by, school_id)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, grade.Name, grade.Description, grade.VisualTheme, grade.CreatedBy, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}

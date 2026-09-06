package grade

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Grade, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description, ''), visual_theme, created_by, created_at
		FROM grades
		WHERE id = $1 AND ($2 = '' OR school_id = NULLIF($2, '')::uuid)
	`, id, tenant.SchoolID(ctx))

	var grade domain.Grade
	if err := row.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.VisualTheme, &grade.CreatedBy, &grade.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

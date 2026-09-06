package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, grade domain.Grade) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE grades SET name = $1, description = $2, visual_theme = $3 WHERE id = $4 AND ($5 = '' OR school_id = NULLIF($5, '')::uuid)
	`, grade.Name, grade.Description, grade.VisualTheme, id, tenant.SchoolID(ctx))
	return err
}

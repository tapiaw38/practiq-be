package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, grade domain.Grade) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE grades SET name = $1, description = $2, visual_theme = $3 WHERE id = $4
	`, grade.Name, grade.Description, grade.VisualTheme, id)
	return err
}

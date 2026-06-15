package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, grade domain.Grade) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE grades SET name = $1, description = $2 WHERE id = $3
	`, grade.Name, grade.Description, id)
	return err
}

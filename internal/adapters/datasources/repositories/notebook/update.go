package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, n domain.Notebook) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebooks SET title = $1, description = $2, updated_at = NOW() WHERE id = $3
	`, n.Title, n.Description, id)
	return err
}

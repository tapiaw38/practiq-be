package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, t domain.Topic) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE topics SET title = $1, description = $2, order_index = $3 WHERE id = $4
	`, t.Title, t.Description, t.OrderIndex, id)
	return err
}

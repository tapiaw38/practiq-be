package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, t domain.Topic) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE topics SET title = $1, description = $2, order_index = $3
		WHERE id = $4 AND ($5 = '' OR school_id = NULLIF($5, '')::uuid)
	`, t.Title, t.Description, t.OrderIndex, id, tenant.SchoolID(ctx))
	return err
}

package exercise

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM exercises WHERE id=$1 AND ($2 = '' OR school_id = NULLIF($2, '')::uuid)`, id, tenant.SchoolID(ctx))
	return err
}

package course

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `UPDATE courses SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL AND ($3 = '' OR school_id = NULLIF($3, '')::uuid)`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id, tenant.SchoolID(ctx))
	return err
}

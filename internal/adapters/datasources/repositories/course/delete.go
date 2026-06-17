package course

import (
	"context"
	"time"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `UPDATE courses SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

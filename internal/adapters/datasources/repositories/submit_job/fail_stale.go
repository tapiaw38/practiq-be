package submitjob

import (
	"context"
	"time"
)

func (r *repository) FailStale(ctx context.Context, olderThan time.Duration) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE submit_jobs
		SET status = 'failed',
		    error_code = 'submit:interrupted',
		    message = 'La entrega se interrumpió. Volvé a enviarla.',
		    updated_at = NOW()
		WHERE status = 'processing' AND updated_at < $1
	`, time.Now().UTC().Add(-olderThan))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

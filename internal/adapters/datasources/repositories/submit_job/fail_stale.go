package submitjob

import (
	"context"
	"time"
)

// FailStale closes jobs that can no longer finish.
//
// A submission runs in a goroutine and the payload lives only there, so a
// process that restarts mid-flight leaves its jobs claiming to be processing
// with nothing behind them. The client polls that status forever: the student
// sees a spinner for work nobody is doing, and does not resubmit because the
// app says it is still going.
//
// Marking them failed is what ends the wait. It does not retry them — retrying
// would mean storing the answer and spending the assistant's calls again on a
// delivery the student may have already replaced — it tells the truth, and
// resubmitting is safe now that submissions are versioned.
//
// olderThan must exceed the longest a live job can legitimately take, or this
// would fail submissions that are still running.
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

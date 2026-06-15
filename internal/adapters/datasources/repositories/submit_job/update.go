package submitjob

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, job domain.SubmitJob) error {
	query := `
		UPDATE submit_jobs
		SET status = $2, error_code = NULLIF($3,''), message = NULLIF($4,''), result = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Status, job.ErrorCode, job.Message, nullableJSON(job.Result), time.Now().UTC(),
	)
	return err
}

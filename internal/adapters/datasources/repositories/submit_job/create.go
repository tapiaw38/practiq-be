package submitjob

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, job domain.SubmitJob) error {
	query := `
		INSERT INTO submit_jobs (id, kind, status, error_code, message, result, created_at, updated_at)
		VALUES ($1, $2, $3, NULLIF($4,''), NULLIF($5,''), $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Kind, job.Status, job.ErrorCode, job.Message, nullableJSON(job.Result), job.CreatedAt, job.UpdatedAt,
	)
	return err
}

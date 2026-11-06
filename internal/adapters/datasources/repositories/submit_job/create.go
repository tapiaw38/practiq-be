package submitjob

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, job domain.SubmitJob) error {
	query := `
		INSERT INTO submit_jobs (id, kind, student_id, status, error_code, message, result, created_at, updated_at)
		VALUES ($1, $2, NULLIF($3,''), $4, NULLIF($5,''), NULLIF($6,''), $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Kind, job.StudentID, job.Status, job.ErrorCode, job.Message, nullableJSON(job.Result), job.CreatedAt, job.UpdatedAt,
	)
	return err
}

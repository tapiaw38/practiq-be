package submitjob

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetByID(ctx context.Context, id string) (*domain.SubmitJob, error) {
	query := `
		SELECT id, kind, status, COALESCE(error_code,''), COALESCE(message,''), result, created_at, updated_at
		FROM submit_jobs
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var job domain.SubmitJob
	err := row.Scan(&job.ID, &job.Kind, &job.Status, &job.ErrorCode, &job.Message, &job.Result, &job.CreatedAt, &job.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

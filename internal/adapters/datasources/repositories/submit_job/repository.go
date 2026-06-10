package submitjob

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, job domain.SubmitJob) error
	Update(ctx context.Context, job domain.SubmitJob) error
	GetByID(ctx context.Context, id string) (*domain.SubmitJob, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

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

func nullableJSON(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	return raw
}

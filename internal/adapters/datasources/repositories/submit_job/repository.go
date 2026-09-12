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
	// FailStale closes jobs left processing by a process that is gone, so a
	// student stops waiting on work nobody is doing.
	FailStale(ctx context.Context, olderThan time.Duration) (int64, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

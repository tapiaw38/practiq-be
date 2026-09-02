package submitjob

import (
	"context"
	"database/sql"

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

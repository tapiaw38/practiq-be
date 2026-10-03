package coursecuriosities

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, courseID string) (*domain.CourseCuriosities, error)
	Upsert(ctx context.Context, c domain.CourseCuriosities) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

package courseprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, studentID, courseID string) (*domain.StudentCourseProgress, error)
	ListByStudent(ctx context.Context, studentID string) ([]domain.StudentCourseProgress, error)
	Upsert(ctx context.Context, studentID, courseID string, level int) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

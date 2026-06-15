package course

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type ListFilterOptions struct {
	TeacherID string
	StudentID string
}
type Repository interface {
	Create(context.Context, domain.Course) (string, error)
	Get(context.Context, string) (*domain.Course, error)
	List(context.Context, ListFilterOptions) ([]domain.Course, error)
	Update(context.Context, string, domain.Course) error
	Delete(context.Context, string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

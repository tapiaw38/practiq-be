package exercise

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Exercise) (string, error)
	Get(context.Context, string) (*domain.Exercise, error)
	List(context.Context, ListFilter) ([]domain.Exercise, error)
	Update(context.Context, string, domain.Exercise) error
	Delete(context.Context, string) error
}
type ListFilter struct {
	TopicID string
	Limit   int
	Offset  int
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

package subject

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Subject) (string, error)
	List(context.Context) ([]domain.Subject, error)
	Get(context.Context, string) (*domain.Subject, error)
	Update(context.Context, string, domain.Subject) error
	Delete(context.Context, string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

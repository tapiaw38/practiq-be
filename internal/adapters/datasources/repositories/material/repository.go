package material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Material) (string, error)
	List(context.Context, string) ([]domain.Material, error)
	Get(context.Context, string) (*domain.Material, error)
	Update(context.Context, string, domain.Material) error
	Delete(context.Context, string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

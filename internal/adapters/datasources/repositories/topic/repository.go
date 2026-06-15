package topic

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Topic) (string, error)
	List(context.Context, string) ([]domain.Topic, error)
	Get(context.Context, string) (*domain.Topic, error)
	Update(context.Context, string, domain.Topic) error
	Delete(context.Context, string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

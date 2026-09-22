package userprofile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.UserProfile) error
	Get(context.Context, string) (*domain.UserProfile, error)
	UpdateUITheme(context.Context, string, string) error
	UpdateAvatarSeed(ctx context.Context, id, seed string) error
	UpdateAcademicStatus(context.Context, string, string) error
	UpdateProfileType(context.Context, string, string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

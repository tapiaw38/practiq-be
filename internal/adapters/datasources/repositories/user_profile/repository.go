package userprofile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.UserProfile) error
	Get(context.Context, string) (*domain.UserProfile, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, p domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, name, email, profile_type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, profile_type = $4
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Email, p.ProfileType)
	return err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.UserProfile, error) {
	query := `SELECT id, name, email, profile_type, created_at FROM user_profiles WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p domain.UserProfile
	err := row.Scan(&p.ID, &p.Name, &p.Email, &p.ProfileType, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

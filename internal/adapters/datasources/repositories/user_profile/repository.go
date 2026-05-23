package userprofile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.UserProfile) error
	Get(context.Context, string) (*domain.UserProfile, error)
	UpdateAssistantConfig(context.Context, string, string, string) error
	UpdateAcademicStatus(context.Context, string, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, p domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, name, email, profile_type, academic_status, assistant_base_url, assistant_api_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			email = $3,
			profile_type = $4,
			assistant_base_url = $6,
			assistant_api_key = $7
	`
	status := p.AcademicStatus
	if status == "" {
		status = "active"
	}
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Email, p.ProfileType, status, p.AssistantBaseURL, p.AssistantAPIKey)
	return err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.UserProfile, error) {
	query := `SELECT id, name, email, profile_type, academic_status, assistant_base_url, assistant_api_key, created_at FROM user_profiles WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p domain.UserProfile
	err := row.Scan(&p.ID, &p.Name, &p.Email, &p.ProfileType, &p.AcademicStatus, &p.AssistantBaseURL, &p.AssistantAPIKey, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) UpdateAssistantConfig(ctx context.Context, id, baseURL, apiKey string) error {
	query := `
		UPDATE user_profiles
		SET assistant_base_url = $2, assistant_api_key = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, baseURL, apiKey)
	return err
}

func (r *repository) UpdateAcademicStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE user_profiles
		SET academic_status = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

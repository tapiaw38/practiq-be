package userprofile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.UserProfile, error) {
	query := `SELECT id, name, email, profile_type, academic_status, COALESCE(timezone,''), assistant_base_url, assistant_api_key, created_at FROM user_profiles WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p domain.UserProfile
	err := row.Scan(&p.ID, &p.Name, &p.Email, &p.ProfileType, &p.AcademicStatus, &p.Timezone, &p.AssistantBaseURL, &p.AssistantAPIKey, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

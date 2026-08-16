package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, p domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, name, email, profile_type, academic_status, assistant_base_url, assistant_api_key, timezone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			email = $3,
			profile_type = $4,
			assistant_base_url = $6,
			assistant_api_key = $7,
			-- Only overwrite when the caller actually supplies one; most
			-- callers do not know the timezone and must not erase it.
			timezone = COALESCE(NULLIF($8, ''), user_profiles.timezone)
	`
	status := p.AcademicStatus
	if status == "" {
		status = "active"
	}
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Email, p.ProfileType, status, p.AssistantBaseURL, p.AssistantAPIKey, p.Timezone)
	return err
}

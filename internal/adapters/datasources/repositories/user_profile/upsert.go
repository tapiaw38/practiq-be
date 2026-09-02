package userprofile

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, p domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, profile_type, academic_status, assistant_base_url, assistant_api_key, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			profile_type = $2,
			assistant_base_url = $4,
			assistant_api_key = $5,
			-- Only overwrite when the caller actually supplies one; most
			-- callers do not know the timezone and must not erase it.
			timezone = COALESCE(NULLIF($6, ''), user_profiles.timezone)
	`
	status := p.AcademicStatus
	if status == "" {
		status = "active"
	}
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ProfileType, status, p.AssistantBaseURL, p.AssistantAPIKey, p.Timezone)
	return err
}

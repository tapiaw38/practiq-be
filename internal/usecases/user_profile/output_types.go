package userprofile

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	ProfileData struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		Email            string `json:"email"`
		ProfileType      string `json:"profile_type"`
		AcademicStatus   string `json:"academic_status"`
		AssistantBaseURL string `json:"assistant_base_url"`
		AssistantAPIKey  string `json:"assistant_api_key"`
		UITheme          string `json:"ui_theme"`
		CreatedAt        string `json:"created_at"`
	}
)

// toProfileData takes name/email pre-resolved by the caller (from
// auth-api-be, see internal/platform/identity) since domain.UserProfile no
// longer carries identity fields.
func toProfileData(p domain.UserProfile, name, email string) ProfileData {
	return ProfileData{
		ID:               p.ID,
		Name:             name,
		Email:            email,
		ProfileType:      p.ProfileType,
		AcademicStatus:   p.AcademicStatus,
		AssistantBaseURL: p.AssistantBaseURL,
		AssistantAPIKey:  p.AssistantAPIKey,
		UITheme:          p.UITheme,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

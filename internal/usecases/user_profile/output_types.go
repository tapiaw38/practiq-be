package userprofile

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	ProfileData struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		ProfileType    string `json:"profile_type"`
		AcademicStatus string `json:"academic_status"`
		UITheme        string `json:"ui_theme"`
		// AssistantEnabled tells the client whether to show the assistant. It
		// replaces the per-profile credentials the client used to receive and
		// gate on, which meant handing the key to every browser.
		AssistantEnabled bool   `json:"assistant_enabled"`
		CreatedAt        string `json:"created_at"`
	}
)

// toProfileData takes name/email pre-resolved by the caller (from
// auth-api-be, see internal/platform/identity) since domain.UserProfile no
// longer carries identity fields.
func toProfileData(p domain.UserProfile, name, email string, assistantEnabled bool) ProfileData {
	return ProfileData{
		ID:               p.ID,
		Name:             name,
		Email:            email,
		ProfileType:      p.ProfileType,
		AcademicStatus:   p.AcademicStatus,
		UITheme:          p.UITheme,
		AssistantEnabled: assistantEnabled,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

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
		AvatarSeed     string `json:"avatar_seed"`

		AssistantEnabled bool   `json:"assistant_enabled"`
		CreatedAt        string `json:"created_at"`
	}
)

func toProfileData(p domain.UserProfile, name, email string, assistantEnabled bool) ProfileData {
	return ProfileData{
		ID:               p.ID,
		Name:             name,
		Email:            email,
		ProfileType:      p.ProfileType,
		AcademicStatus:   p.AcademicStatus,
		UITheme:          p.UITheme,
		AvatarSeed:       p.AvatarSeed,
		AssistantEnabled: assistantEnabled,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

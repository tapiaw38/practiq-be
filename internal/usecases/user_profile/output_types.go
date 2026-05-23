package userprofile

import "github.com/tapiaw38/practiq-be/internal/domain"

type ProfileData struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	ProfileType      string `json:"profile_type"`
	AcademicStatus   string `json:"academic_status"`
	AssistantBaseURL string `json:"assistant_base_url"`
	AssistantAPIKey  string `json:"assistant_api_key"`
	CreatedAt        string `json:"created_at"`
}

type ProfileOutput struct {
	Data ProfileData `json:"data"`
}

func toProfileData(p domain.UserProfile) ProfileData {
	return ProfileData{
		ID:               p.ID,
		Name:             p.Name,
		Email:            p.Email,
		ProfileType:      p.ProfileType,
		AcademicStatus:   p.AcademicStatus,
		AssistantBaseURL: p.AssistantBaseURL,
		AssistantAPIKey:  p.AssistantAPIKey,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

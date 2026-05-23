package domain

import "time"

type UserProfile struct {
	ID               string
	Name             string
	Email            string
	ProfileType      string
	AcademicStatus   string
	AssistantBaseURL string
	AssistantAPIKey  string
	CreatedAt        time.Time
}

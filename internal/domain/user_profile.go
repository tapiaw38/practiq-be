package domain

import "time"

type UserProfile struct {
	ID             string
	Name           string
	Email          string
	ProfileType    string
	AcademicStatus string
	// Timezone is the IANA zone the student's day is measured in, reported by
	// the browser. Empty falls back to DefaultTimezone.
	Timezone         string
	AssistantBaseURL string
	AssistantAPIKey  string
	UITheme          string
	CreatedAt        time.Time
}

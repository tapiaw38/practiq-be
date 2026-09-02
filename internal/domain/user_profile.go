package domain

import "time"

// UserProfile deliberately carries no identity fields — auth-api-be is the
// single source of truth for name/email. Callers that need display name
// resolve it via internal/platform/identity instead of a local cache.
type UserProfile struct {
	ID             string
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

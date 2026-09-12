package domain

import "time"

type Grade struct {
	ID string
	// SchoolID is who owns it. A grade or subject belongs to exactly one
	// school; a second school needing the same name gets its own row.
	SchoolID    string
	Name        string
	Description string
	VisualTheme string
	CreatedBy   string
	CreatedAt   time.Time
}

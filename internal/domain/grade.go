package domain

import "time"

type Grade struct {
	ID string

	SchoolID    string
	Name        string
	Description string
	VisualTheme string
	CreatedBy   string
	CreatedAt   time.Time
}

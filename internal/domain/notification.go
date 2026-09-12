package domain

import "time"

// Notification types.
const (
	NotificationLevelTestScheduled = "level_test_scheduled"
)

// Notification resource types.
const (
	NotificationResourcePracticeSheet = "practice_sheet"
)

type Notification struct {
	ID           string
	UserID       string
	Type         string
	Title        string
	Body         string
	ResourceType string
	ResourceID   string
	// ScheduledAt is when the referenced event happens, nil when it has no date.
	ScheduledAt *time.Time
	ReadAt      *time.Time
	CreatedAt   time.Time
}

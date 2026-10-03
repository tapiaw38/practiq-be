package domain

import "time"

const (
	NotificationLevelTestScheduled = "level_test_scheduled"
)

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

	ScheduledAt *time.Time
	ReadAt      *time.Time
	CreatedAt   time.Time
}

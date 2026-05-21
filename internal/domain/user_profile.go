package domain

import "time"

type UserProfile struct {
	ID          string
	Name        string
	Email       string
	ProfileType string
	CreatedAt   time.Time
}

package domain

import "time"

type Subject struct {
	ID string

	SchoolID    string
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
}

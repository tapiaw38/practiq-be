package domain

import "time"

type Subject struct {
	ID          string
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
}

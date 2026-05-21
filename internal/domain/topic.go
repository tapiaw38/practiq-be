package domain

import "time"

type Topic struct {
	ID          string
	CourseID    string
	Title       string
	Description string
	OrderIndex  int
	CreatedAt   time.Time
}

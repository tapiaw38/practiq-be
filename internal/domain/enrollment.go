package domain

import "time"

type Enrollment struct {
	ID        string
	CourseID  string
	StudentID string
	Status    string
	CreatedAt time.Time
}

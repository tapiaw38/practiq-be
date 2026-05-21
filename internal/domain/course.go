package domain

import "time"

type Course struct {
	ID          string
	TeacherID   string
	Title       string
	Description string
	Level       string
	Subject     string
	CreatedAt   time.Time
}

package domain

import "time"

type StudentCourseProgress struct {
	ID           string
	StudentID    string
	CourseID     string
	CurrentLevel int
	UpdatedAt    time.Time
}

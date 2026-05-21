package domain

import "time"

type Material struct {
	ID            string
	CourseID      string
	TeacherID     string
	Title         string
	Type          string
	FileURL       string
	ExtractedText string
	Status        string
	CreatedAt     time.Time
}

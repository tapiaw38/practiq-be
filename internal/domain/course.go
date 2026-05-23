package domain

import "time"

type Course struct {
	ID          string
	TeacherID   string
	GradeID     string
	GradeName   string
	SubjectID   string
	SubjectName string
	Title       string
	Description string
	Level       string
	Subject     string
	CreatedAt   time.Time
}

package domain

import "time"

type TeacherStudentAssignment struct {
	ID        string
	TeacherID string
	StudentID string
	Status    string
	CreatedAt time.Time
}

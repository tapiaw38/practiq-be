package domain

import "time"

type GradeMembership struct {
	ID        string
	GradeID   string
	UserID    string
	CreatedAt time.Time
}

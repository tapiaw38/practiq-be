package studentcoursexp

import (
	"context"
	"database/sql"
)

type AwardInput struct {
	StudentID       string
	CourseID        string
	EventKey        string
	EventType       string
	Points          int
	PracticeSheetID string
	ExerciseID      string
}

type AwardResult struct {
	Awarded bool
	TotalXP int
}

type CourseXP struct {
	CourseID string
	TotalXP  int
}

type Repository interface {
	Award(context.Context, AwardInput) (AwardResult, error)
	Balance(ctx context.Context, studentID, courseID string) (int, error)
	ListByStudent(context.Context, string) ([]CourseXP, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

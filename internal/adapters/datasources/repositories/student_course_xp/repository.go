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

type LeaderboardEntry struct {
	StudentID string
	// Full name as stored. Shortening it for display is the usecase's job.
	Name     string
	TotalXP  int
	Position int
}

type Repository interface {
	Award(context.Context, AwardInput) (AwardResult, error)
	Balance(ctx context.Context, studentID, courseID string) (int, error)
	LeaderboardByCourse(ctx context.Context, courseID, studentID string, top int) ([]LeaderboardEntry, error)
	ListByStudent(context.Context, string) ([]CourseXP, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

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

// LeaderboardEntry carries no name: identity lives in auth-api-be, not in
// this database. The usecase resolves it.
type LeaderboardEntry struct {
	StudentID string
	TotalXP   int
	Position  int
}

type Repository interface {
	Award(context.Context, AwardInput) (AwardResult, error)
	Balance(ctx context.Context, studentID, courseID string) (int, error)
	LeaderboardByCourse(ctx context.Context, courseID, studentID string, top int) ([]LeaderboardEntry, error)
	ListByStudent(context.Context, string) ([]CourseXP, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

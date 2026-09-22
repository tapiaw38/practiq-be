package studentpracticestate

import (
	"context"
	"database/sql"
)

type Repository interface {
	MarkOpened(ctx context.Context, studentID, sheetID string) error
	GetLastOpenedPractice(ctx context.Context, studentID string) (*ResumePractice, error)
}

type ResumePractice struct {
	SheetID    string
	TopicID    string
	TopicTitle string
	Level      int
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

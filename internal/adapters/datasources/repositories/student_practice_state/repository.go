package studentpracticestate

import (
	"context"
	"database/sql"
)

type Repository interface {
	MarkOpened(ctx context.Context, studentID, sheetID string) error
	GetLastOpenedSheetID(ctx context.Context, studentID string) (string, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

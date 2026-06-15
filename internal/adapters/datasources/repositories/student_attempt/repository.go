package studentattempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.StudentAttempt) (string, error)
	ListBySheet(ctx context.Context, studentID, sheetID string) ([]domain.StudentAttempt, error)
	SaveCanvasWork(ctx context.Context, attemptID, imageData string) error
	GetLastPracticedSheetID(ctx context.Context, studentID string) (string, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

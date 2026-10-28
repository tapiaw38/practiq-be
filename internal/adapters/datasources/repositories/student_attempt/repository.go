package studentattempt

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.StudentAttempt) (string, error)

	ClaimLevelTestSubmission(ctx context.Context, studentID, sheetID string, maxAttempts int) (bool, error)

	LevelTestProgress(ctx context.Context, studentID, sheetID string) (int, *time.Time, error)

	MarkLevelTestStarted(ctx context.Context, studentID, sheetID string) error

	CloseExpiredLevelTest(ctx context.Context, studentID, sheetID string, deadline time.Time) error

	ReleaseLevelTestSubmission(ctx context.Context, studentID, sheetID string) error

	DeleteBySheet(ctx context.Context, studentID, sheetID string) error
	ListBySheet(ctx context.Context, studentID, sheetID string) ([]domain.StudentAttempt, error)
	SaveCanvasWork(ctx context.Context, attemptID, imageData string) error
	GetLastPracticedSheetID(ctx context.Context, studentID string) (string, error)
	GetDailyAttempts(ctx context.Context, studentID, courseID string, from, to *time.Time) ([]domain.DailyAttemptCount, error)
	ListPendingReview(ctx context.Context, filter PendingReviewFilter) ([]domain.PendingAttemptReview, error)
	GetTeacherForAttempt(ctx context.Context, attemptID string) (string, error)
	GetExerciseIDForAttempt(ctx context.Context, attemptID string) (string, error)
	Review(ctx context.Context, attemptID string, isCorrect bool, feedback string) error
	GetSheetOutcome(ctx context.Context, studentID, sheetID string) (domain.SheetOutcome, error)
	GetAttemptContext(ctx context.Context, attemptID string) (domain.AttemptContext, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

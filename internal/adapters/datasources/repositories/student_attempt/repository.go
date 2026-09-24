package studentattempt

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.StudentAttempt) (string, error)
	// ClaimLevelTestSubmission atomically reserves a student's one submission
	// for a level test. False means it was already submitted.
	ClaimLevelTestSubmission(ctx context.Context, studentID, sheetID string, maxAttempts int) (bool, error)
	// LevelTestProgress reports how many times the student has submitted this
	// test and when they first opened it. A student who never opened it has
	// zero attempts and no start.
	LevelTestProgress(ctx context.Context, studentID, sheetID string) (int, *time.Time, error)
	// MarkLevelTestStarted records when the student first opened the test,
	// which is when a time limit starts running. Calling it again keeps the
	// original moment: reopening the page must not buy more time.
	MarkLevelTestStarted(ctx context.Context, studentID, sheetID string) error
	// CloseExpiredLevelTest clears an active test window only when it still
	// belongs to the expired attempt. The conditional update makes concurrent
	// expiry requests safe and cannot close a later attempt.
	CloseExpiredLevelTest(ctx context.Context, studentID, sheetID string, deadline time.Time) error
	// ReleaseLevelTestSubmission undoes a claim whose submission never landed.
	ReleaseLevelTestSubmission(ctx context.Context, studentID, sheetID string) error
	// DeleteBySheet removes an incomplete level-test submission before releasing
	// its claim. It is only safe after the claim has established exclusivity.
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

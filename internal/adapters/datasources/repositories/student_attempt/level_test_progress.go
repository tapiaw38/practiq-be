package studentattempt

import (
	"context"
	"database/sql"
	"time"
)

func (r *repository) LevelTestProgress(ctx context.Context, studentID, sheetID string) (int, *time.Time, error) {
	var (
		attempts  int
		startedAt *time.Time
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT attempts, started_at
		FROM level_test_submissions
		WHERE student_id = $1 AND practice_sheet_id = $2::uuid
	`, studentID, sheetID).Scan(&attempts, &startedAt)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	return attempts, startedAt, err
}

// MarkLevelTestStarted creates the row with no attempt spent, so the clock can
// start before anything is submitted. COALESCE keeps the first moment: a
// student who reloads the page, or opens it on a second device, does not get
// the time limit reset.
func (r *repository) MarkLevelTestStarted(ctx context.Context, studentID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO level_test_submissions (practice_sheet_id, student_id, attempts, started_at, submitted_at)
		VALUES ($1::uuid, $2, 0, NOW(), NULL)
		ON CONFLICT (practice_sheet_id, student_id) DO UPDATE
		SET started_at = COALESCE(level_test_submissions.started_at, NOW())
	`, sheetID, studentID)
	return err
}

// CloseExpiredLevelTest closes only the window that has reached deadline.
// A student may have reopened the test while an earlier expiry request was in
// flight; matching the deadline prevents that request from clearing the new
// window.
func (r *repository) CloseExpiredLevelTest(ctx context.Context, studentID, sheetID string, deadline time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE level_test_submissions
		SET started_at = NULL
		WHERE student_id = $1
		  AND practice_sheet_id = $2::uuid
		  AND started_at IS NOT NULL
		  AND started_at <= $3
	`, studentID, sheetID, deadline)
	return err
}

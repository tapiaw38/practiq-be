package studentattempt

import "context"

// ReleaseLevelTestSubmission undoes a claim whose submission never persisted.
//
// The claim is committed before the attempts are written, and the repositories
// share a *sql.DB with no transaction plumbing. Without this, a transient
// failure after the claim left the student permanently locked out: every retry
// answered "already submitted" while no attempt and no promotion existed.
func (r *repository) ReleaseLevelTestSubmission(ctx context.Context, studentID, sheetID string) error {
	// Give the attempt back rather than dropping the row: it also holds when
	// the student started, and deleting it would hand them a fresh clock.
	_, err := r.db.ExecContext(ctx, `
		UPDATE level_test_submissions
		SET attempts = GREATEST(attempts - 1, 0)
		WHERE practice_sheet_id = $1::uuid AND student_id = $2
	`, sheetID, studentID)
	return err
}

package studentattempt

import "context"

// ReleaseLevelTestSubmission undoes a claim whose submission never persisted.
//
// The claim is committed before the attempts are written, and the repositories
// share a *sql.DB with no transaction plumbing. Without this, a transient
// failure after the claim left the student permanently locked out: every retry
// answered "already submitted" while no attempt and no promotion existed.
func (r *repository) ReleaseLevelTestSubmission(ctx context.Context, studentID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM level_test_submissions
		WHERE practice_sheet_id = $1::uuid AND student_id = $2
	`, sheetID, studentID)
	return err
}

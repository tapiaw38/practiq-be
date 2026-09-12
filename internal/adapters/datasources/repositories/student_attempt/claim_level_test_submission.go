package studentattempt

import "context"

// ClaimLevelTestSubmission is intentionally an INSERT with a unique key, not
// a read-then-write check. Async submit jobs can run concurrently, and both
// would otherwise observe an empty attempt list before either writes answers.
func (r *repository) ClaimLevelTestSubmission(ctx context.Context, studentID, sheetID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO level_test_submissions (practice_sheet_id, student_id)
		VALUES ($1::uuid, $2)
		ON CONFLICT (practice_sheet_id, student_id) DO NOTHING
	`, sheetID, studentID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected == 1, nil
}

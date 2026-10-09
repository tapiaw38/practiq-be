package studentattempt

import "context"

func (r *repository) ClaimLevelTestSubmission(ctx context.Context, studentID, sheetID string, maxAttempts int) (bool, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO level_test_submissions (practice_sheet_id, student_id, attempts)
		VALUES ($1::uuid, $2, 1)
		ON CONFLICT (practice_sheet_id, student_id) DO UPDATE
		SET attempts = level_test_submissions.attempts + 1,
		    submitted_at = NOW(),
		    started_at = NULL
		WHERE level_test_submissions.attempts < $3
	`, sheetID, studentID, maxAttempts)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected == 1, nil
}

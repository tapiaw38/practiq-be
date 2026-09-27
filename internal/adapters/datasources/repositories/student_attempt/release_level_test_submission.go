package studentattempt

import "context"

func (r *repository) ReleaseLevelTestSubmission(ctx context.Context, studentID, sheetID string) error {

	_, err := r.db.ExecContext(ctx, `
		UPDATE level_test_submissions
		SET attempts = GREATEST(attempts - 1, 0)
		WHERE practice_sheet_id = $1::uuid AND student_id = $2
	`, sheetID, studentID)
	return err
}

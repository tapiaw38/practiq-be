package studentattempt

import "context"

func (r *repository) HasLevelTestSubmission(ctx context.Context, studentID, sheetID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM level_test_submissions
			WHERE student_id = $1 AND practice_sheet_id = $2::uuid
		)
	`, studentID, sheetID).Scan(&exists)
	return exists, err
}

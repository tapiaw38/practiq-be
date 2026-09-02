package studentattempt

import "context"

func (r *repository) DeleteBySheet(ctx context.Context, studentID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM student_attempts
		WHERE student_id = $1 AND practice_sheet_id = $2::uuid
	`, studentID, sheetID)
	return err
}

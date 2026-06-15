package studentattempt

import (
	"context"
	"database/sql"
)

func (r *repository) GetLastPracticedSheetID(ctx context.Context, studentID string) (string, error) {
	query := `
		SELECT COALESCE(practice_sheet_id::text, '')
		FROM student_attempts
		WHERE student_id = $1 AND practice_sheet_id IS NOT NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	var sheetID string
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&sheetID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return sheetID, err
}

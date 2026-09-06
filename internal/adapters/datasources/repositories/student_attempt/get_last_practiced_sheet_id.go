package studentattempt

import (
	"context"
	"database/sql"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) GetLastPracticedSheetID(ctx context.Context, studentID string) (string, error) {
	// Joined through an active course: PracticeSheet.Get rejects sheets whose
	// course is soft-deleted, so without this the progress payload advertised a
	// "continue where you left off" sheet that could not be opened.
	query := `
		SELECT COALESCE(sa.practice_sheet_id::text, '')
		FROM student_attempts sa
		JOIN practice_sheets ps ON ps.id = sa.practice_sheet_id
		JOIN courses c ON c.id = ps.course_id AND c.deleted_at IS NULL
		WHERE sa.student_id = $1 AND sa.practice_sheet_id IS NOT NULL
		  AND ($2 = '' OR sa.school_id = NULLIF($2, '')::uuid)
		ORDER BY sa.created_at DESC
		LIMIT 1
	`
	var sheetID string
	err := r.db.QueryRowContext(ctx, query, studentID, tenant.SchoolID(ctx)).Scan(&sheetID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return sheetID, err
}

package studentpracticestate

import (
	"context"
	"database/sql"
)

func (r *repository) MarkOpened(ctx context.Context, studentID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_practice_state (student_id, practice_sheet_id, last_opened_at, updated_at)
		VALUES ($1, $2::uuid, now(), now())
		ON CONFLICT (student_id) DO UPDATE
		SET practice_sheet_id = EXCLUDED.practice_sheet_id,
			last_opened_at = EXCLUDED.last_opened_at,
			updated_at = EXCLUDED.updated_at`, studentID, sheetID)
	return err
}

func (r *repository) GetLastOpenedSheetID(ctx context.Context, studentID string) (string, error) {
	// Match the student's course visibility rule. A stale state row must never
	// advertise a draft, deleted or no-longer-accessible sheet.
	var sheetID string
	err := r.db.QueryRowContext(ctx, `
		SELECT s.practice_sheet_id::text
		FROM student_practice_state s
		JOIN practice_sheets ps ON ps.id = s.practice_sheet_id AND ps.sheet_type = 'practice'
		JOIN courses c ON c.id = ps.course_id
		WHERE s.student_id = $1
		  AND s.practice_sheet_id IS NOT NULL
		  AND c.deleted_at IS NULL
		  AND c.status IN ('published', 'archived')
		  AND (
			EXISTS (SELECT 1 FROM enrollments e WHERE e.course_id = c.id AND e.student_id = s.student_id)
			OR EXISTS (SELECT 1 FROM grade_memberships gm WHERE gm.grade_id = c.grade_id AND gm.user_id = s.student_id)
		  )
		LIMIT 1`, studentID).Scan(&sheetID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return sheetID, err
}

package practicesheet

import "context"

func (r *repository) HasOtherLevelTest(ctx context.Context, courseID string, level int, excludeID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM practice_sheets
			WHERE course_id = $1
			  AND level = $2
			  AND sheet_type = 'level_test'
			  AND ($3 = '' OR id <> $3::uuid)
		)
	`, courseID, level, excludeID).Scan(&exists)
	return exists, err
}

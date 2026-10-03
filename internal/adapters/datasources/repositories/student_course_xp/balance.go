package studentcoursexp

import (
	"context"
	"database/sql"
	"errors"
)

func (r *repository) Balance(ctx context.Context, studentID, courseID string) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `
		SELECT total_xp FROM student_course_xp
		WHERE student_id = $1 AND course_id = $2::uuid`, studentID, courseID,
	).Scan(&total)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return total, nil
}

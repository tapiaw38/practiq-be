package courseprogress

import "context"

func (r *repository) Upsert(ctx context.Context, studentID, courseID string, level int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_course_progress (student_id, course_id, current_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (student_id, course_id)
		DO UPDATE SET current_level = GREATEST(student_course_progress.current_level, EXCLUDED.current_level),
		              updated_at = NOW()
	`, studentID, courseID, level)
	return err
}

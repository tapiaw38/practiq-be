package courseprogress

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Upsert(ctx context.Context, studentID, courseID string, level int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_course_progress (student_id, course_id, current_level, school_id)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid)
		ON CONFLICT (student_id, course_id)
		DO UPDATE SET current_level = GREATEST(student_course_progress.current_level, EXCLUDED.current_level),
		              school_id = COALESCE(student_course_progress.school_id, EXCLUDED.school_id),
		              updated_at = NOW()
	`, studentID, courseID, level, tenant.SchoolID(ctx))
	return err
}

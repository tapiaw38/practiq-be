package courseprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, studentID, courseID string) (*domain.StudentCourseProgress, error) {
	var p domain.StudentCourseProgress
	err := r.db.QueryRowContext(ctx, `
		SELECT id, student_id, course_id, current_level, updated_at
		FROM student_course_progress
		WHERE student_id = $1 AND course_id = $2 AND ($3 = '' OR school_id = NULLIF($3, '')::uuid)
	`, studentID, courseID, tenant.SchoolID(ctx)).Scan(&p.ID, &p.StudentID, &p.CourseID, &p.CurrentLevel, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return &domain.StudentCourseProgress{StudentID: studentID, CourseID: courseID, CurrentLevel: 1}, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

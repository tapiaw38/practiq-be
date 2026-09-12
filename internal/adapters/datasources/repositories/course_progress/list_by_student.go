package courseprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListByStudent(ctx context.Context, studentID string) ([]domain.StudentCourseProgress, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, student_id, course_id, current_level, updated_at
		FROM student_course_progress
		WHERE student_id = $1
		ORDER BY updated_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentCourseProgress
	for rows.Next() {
		var p domain.StudentCourseProgress
		if err := rows.Scan(&p.ID, &p.StudentID, &p.CourseID, &p.CurrentLevel, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

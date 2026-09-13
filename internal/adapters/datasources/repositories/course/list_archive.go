package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// ListArchive deliberately bypasses the active-school filter used by normal
// course reads. It is only exposed through the superadmin archive endpoint.
func (r *repository) ListArchive(ctx context.Context, schoolID string) ([]domain.Course, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.teacher_id, COALESCE(g.school_id::text, s.school_id::text, ''),
		       COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(g.visual_theme, 'primary'),
		       COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title,
		       c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
		WHERE COALESCE(g.school_id, s.school_id)::text = $1
		ORDER BY c.created_at DESC
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []domain.Course{}
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.TeacherID, &c.SchoolID, &c.GradeID, &c.GradeName, &c.GradeTheme, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

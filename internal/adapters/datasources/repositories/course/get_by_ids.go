package course

import (
	"context"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetByIDs(ctx context.Context, ids []string) ([]domain.Course, error) {
	if len(ids) == 0 {
		return []domain.Course{}, nil
	}

	query := `
		SELECT c.id, c.teacher_id, COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(g.visual_theme, 'primary'), COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title, c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
		WHERE c.id = ANY($1) AND c.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]domain.Course, 0, len(ids))
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.TeacherID, &c.GradeID, &c.GradeName, &c.GradeTheme, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, nil
}

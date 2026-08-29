package course

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Course, error) {
	query := `
		SELECT c.id, c.teacher_id, COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(g.visual_theme, 'primary'), COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title, c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var c domain.Course
	err := row.Scan(&c.ID, &c.TeacherID, &c.GradeID, &c.GradeName, &c.GradeTheme, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

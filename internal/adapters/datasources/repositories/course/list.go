package course

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) List(ctx context.Context, opts ListFilterOptions) ([]domain.Course, error) {
	query := `
		SELECT c.id, c.teacher_id, COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(g.visual_theme, 'primary'), COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title, c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
		WHERE c.deleted_at IS NULL
	`
	args := []interface{}{}
	argIdx := 1
	if schoolID := tenant.SchoolID(ctx); schoolID != "" {
		query += fmt.Sprintf(` AND c.school_id = $%d`, argIdx)
		args = append(args, schoolID)
		argIdx++
	}

	if opts.TeacherID != "" {
		query += fmt.Sprintf(` AND c.teacher_id = $%d`, argIdx)
		args = append(args, opts.TeacherID)
		argIdx++
	} else if opts.StudentID != "" {
		query += fmt.Sprintf(` AND (EXISTS (SELECT 1 FROM enrollments e WHERE e.course_id = c.id AND e.student_id = $%d) OR EXISTS (SELECT 1 FROM grade_memberships gm WHERE gm.grade_id = c.grade_id AND gm.user_id = $%d))`, argIdx, argIdx)
		args = append(args, opts.StudentID)
		argIdx++
	}
	_ = argIdx

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.TeacherID, &c.GradeID, &c.GradeName, &c.GradeTheme, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

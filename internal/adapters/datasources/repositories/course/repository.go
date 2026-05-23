package course

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type ListFilterOptions struct {
	TeacherID string
	StudentID string
}

type Repository interface {
	Create(context.Context, domain.Course) (string, error)
	Get(context.Context, string) (*domain.Course, error)
	List(context.Context, ListFilterOptions) ([]domain.Course, error)
	Update(context.Context, string, domain.Course) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c domain.Course) (string, error) {
	query := `
		INSERT INTO courses (teacher_id, grade_id, subject_id, title, description, level, subject)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(
		ctx,
		query,
		c.TeacherID,
		nullableUUID(c.GradeID),
		nullableUUID(c.SubjectID),
		c.Title,
		c.Description,
		c.Level,
		c.Subject,
	).Scan(&id)
	return id, err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Course, error) {
	query := `
		SELECT c.id, c.teacher_id, COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title, c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
		WHERE c.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var c domain.Course
	err := row.Scan(&c.ID, &c.TeacherID, &c.GradeID, &c.GradeName, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(ctx context.Context, opts ListFilterOptions) ([]domain.Course, error) {
	query := `
		SELECT c.id, c.teacher_id, COALESCE(c.grade_id::text, ''), COALESCE(g.name, ''), COALESCE(c.subject_id::text, ''), COALESCE(s.name, c.subject, ''), c.title, c.description, c.level, COALESCE(c.subject, ''), c.created_at
		FROM courses c
		LEFT JOIN grades g ON g.id = c.grade_id
		LEFT JOIN subjects s ON s.id = c.subject_id
	`
	args := []interface{}{}
	argIdx := 1

	if opts.TeacherID != "" {
		query += fmt.Sprintf(` WHERE c.teacher_id = $%d`, argIdx)
		args = append(args, opts.TeacherID)
		argIdx++
	} else if opts.StudentID != "" {
		query += fmt.Sprintf(` WHERE EXISTS (SELECT 1 FROM enrollments e WHERE e.course_id = c.id AND e.student_id = $%d) OR EXISTS (SELECT 1 FROM grade_memberships gm WHERE gm.grade_id = c.grade_id AND gm.user_id = $%d)`, argIdx, argIdx)
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
		if err := rows.Scan(&c.ID, &c.TeacherID, &c.GradeID, &c.GradeName, &c.SubjectID, &c.SubjectName, &c.Title, &c.Description, &c.Level, &c.Subject, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *repository) Update(ctx context.Context, id string, c domain.Course) error {
	query := `UPDATE courses SET title=$1, description=$2, level=$3, subject=$4, grade_id=$5, subject_id=$6 WHERE id=$7`
	_, err := r.db.ExecContext(ctx, query, c.Title, c.Description, c.Level, c.Subject, nullableUUID(c.GradeID), nullableUUID(c.SubjectID), id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return err
}

func nullableUUID(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

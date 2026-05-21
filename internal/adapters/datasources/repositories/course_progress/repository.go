package courseprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, studentID, courseID string) (*domain.StudentCourseProgress, error)
	Upsert(ctx context.Context, studentID, courseID string, level int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Get(ctx context.Context, studentID, courseID string) (*domain.StudentCourseProgress, error) {
	var p domain.StudentCourseProgress
	err := r.db.QueryRowContext(ctx, `
		SELECT id, student_id, course_id, current_level, updated_at
		FROM student_course_progress
		WHERE student_id = $1 AND course_id = $2
	`, studentID, courseID).Scan(&p.ID, &p.StudentID, &p.CourseID, &p.CurrentLevel, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return &domain.StudentCourseProgress{StudentID: studentID, CourseID: courseID, CurrentLevel: 1}, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

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

package enrollment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Enrollment) error
	ListStudents(context.Context, string) ([]domain.UserProfile, error)
	Exists(context.Context, string, string) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, e domain.Enrollment) error {
	query := `INSERT INTO enrollments (course_id, student_id, status) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, e.CourseID, e.StudentID, e.Status)
	return err
}

func (r *repository) ListStudents(ctx context.Context, courseID string) ([]domain.UserProfile, error) {
	query := `
		SELECT up.id, up.name, up.email, up.profile_type, up.created_at
		FROM user_profiles up
		JOIN enrollments e ON e.student_id = up.id
		WHERE e.course_id = $1
		ORDER BY up.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []domain.UserProfile
	for rows.Next() {
		var s domain.UserProfile
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.ProfileType, &s.CreatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *repository) Exists(ctx context.Context, courseID, studentID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM enrollments WHERE course_id=$1 AND student_id=$2`, courseID, studentID).Scan(&count)
	return count > 0, err
}

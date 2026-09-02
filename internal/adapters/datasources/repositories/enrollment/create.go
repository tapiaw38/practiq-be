package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, e domain.Enrollment) error {
	query := `INSERT INTO enrollments (course_id, student_id, status) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, e.CourseID, e.StudentID, e.Status)
	return err
}

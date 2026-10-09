package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

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

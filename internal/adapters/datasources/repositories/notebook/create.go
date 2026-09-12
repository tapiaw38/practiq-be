package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, n domain.Notebook) (string, error) {
	level := n.Level
	if level < 1 {
		level = 1
	}
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebooks (course_id, teacher_id, title, description, level)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, n.CourseID, n.TeacherID, n.Title, n.Description, level).Scan(&id)
	return id, err
}

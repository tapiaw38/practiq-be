package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Create(ctx context.Context, n domain.Notebook) (string, error) {
	level := n.Level
	if level < 1 {
		level = 1
	}
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebooks (course_id, teacher_id, title, description, level, school_id)
		SELECT $1, $2, $3, $4, $5, NULLIF($6,'')::uuid
		WHERE EXISTS (SELECT 1 FROM courses WHERE id = $1 AND ($6 = '' OR school_id = NULLIF($6,'')::uuid))
		RETURNING id
	`, n.CourseID, n.TeacherID, n.Title, n.Description, level, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}

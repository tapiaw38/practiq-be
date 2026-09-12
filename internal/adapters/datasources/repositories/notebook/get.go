package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Notebook, error) {
	var nb domain.Notebook
	err := r.db.QueryRowContext(ctx, `
		SELECT n.id, n.course_id, n.teacher_id, n.title, n.description, n.level, n.created_at, n.updated_at
		FROM notebooks n
		JOIN courses c ON c.id = n.course_id
		WHERE n.id = $1 AND n.deleted_at IS NULL AND c.deleted_at IS NULL
	`, id).Scan(&nb.ID, &nb.CourseID, &nb.TeacherID, &nb.Title, &nb.Description, &nb.Level, &nb.CreatedAt, &nb.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	pages, err := r.listPages(ctx, id)
	if err != nil {
		return nil, err
	}
	nb.Pages = pages
	return &nb, nil
}

package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Notebook, error) {
	var nb domain.Notebook
	err := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, teacher_id, title, description, level, created_at, updated_at
		FROM notebooks WHERE id = $1
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

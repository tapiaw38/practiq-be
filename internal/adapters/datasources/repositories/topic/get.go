package topic

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Topic, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, title, description, order_index, created_at
		FROM topics WHERE id = $1
	`, id)

	var t domain.Topic
	if err := row.Scan(&t.ID, &t.CourseID, &t.Title, &t.Description, &t.OrderIndex, &t.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

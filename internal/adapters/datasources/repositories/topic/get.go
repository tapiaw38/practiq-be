package topic

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Topic, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.course_id, t.title, t.description, t.order_index, t.created_at
		FROM topics t
		JOIN courses c ON c.id = t.course_id
		WHERE t.id = $1 AND c.deleted_at IS NULL
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

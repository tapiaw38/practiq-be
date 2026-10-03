package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, t domain.Topic) (string, error) {
	query := `
		INSERT INTO topics (course_id, title, description, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, t.CourseID, t.Title, t.Description, t.OrderIndex).Scan(&id)
	return id, err
}

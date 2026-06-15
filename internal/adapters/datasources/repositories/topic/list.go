package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Topic, error) {
	query := `SELECT id, course_id, title, description, order_index, created_at FROM topics WHERE course_id = $1 ORDER BY order_index ASC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(&t.ID, &t.CourseID, &t.Title, &t.Description, &t.OrderIndex, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

package topic

import (
	"context"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetByIDs(ctx context.Context, ids []string) ([]domain.Topic, error) {
	if len(ids) == 0 {
		return []domain.Topic{}, nil
	}

	query := `
		SELECT id, course_id, title, description, order_index, created_at
		FROM topics
		WHERE id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := make([]domain.Topic, 0, len(ids))
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(&t.ID, &t.CourseID, &t.Title, &t.Description, &t.OrderIndex, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

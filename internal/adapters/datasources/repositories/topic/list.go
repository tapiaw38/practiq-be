package topic

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Topic, error) {
	query := `
		SELECT t.id, t.course_id, t.title, t.description, t.order_index, t.created_at
		FROM topics t
		JOIN courses c ON c.id = t.course_id
		WHERE t.course_id = $1 AND c.deleted_at IS NULL AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
		ORDER BY t.order_index ASC`

	args := []interface{}{filter.CourseID, tenant.SchoolID(ctx)}
	argIndex := 3

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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

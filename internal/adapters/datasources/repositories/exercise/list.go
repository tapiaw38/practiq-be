package exercise

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Exercise, error) {
	query := `SELECT id, topic_id, COALESCE(material_id::text,''), type, question, COALESCE(correct_answer,''), COALESCE(explanation,''), difficulty, metadata::text, created_at FROM exercises WHERE topic_id = $1 ORDER BY created_at ASC`

	args := []interface{}{filter.TopicID}
	argIndex := 2

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

	var exercises []domain.Exercise
	for rows.Next() {
		var e domain.Exercise
		if err := rows.Scan(&e.ID, &e.TopicID, &e.MaterialID, &e.Type, &e.Question, &e.CorrectAnswer, &e.Explanation, &e.Difficulty, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

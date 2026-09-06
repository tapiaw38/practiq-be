package exercise

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Exercise, error) {
	query := `
		SELECT e.id, e.topic_id, COALESCE(e.material_id::text,''), e.type, e.question, COALESCE(e.correct_answer,''), COALESCE(e.explanation,''), e.difficulty, e.metadata::text, e.created_at
		FROM exercises e
		JOIN topics t ON t.id = e.topic_id
		JOIN courses c ON c.id = t.course_id
		WHERE e.topic_id = $1 AND c.deleted_at IS NULL AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
		ORDER BY e.created_at ASC`

	args := []interface{}{filter.TopicID, tenant.SchoolID(ctx)}
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

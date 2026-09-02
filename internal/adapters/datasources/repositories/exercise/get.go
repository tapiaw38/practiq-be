package exercise

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Exercise, error) {
	query := `SELECT id, topic_id, COALESCE(material_id::text,''), type, question, COALESCE(correct_answer,''), COALESCE(explanation,''), difficulty, metadata::text, created_at FROM exercises WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var e domain.Exercise
	err := row.Scan(&e.ID, &e.TopicID, &e.MaterialID, &e.Type, &e.Question, &e.CorrectAnswer, &e.Explanation, &e.Difficulty, &e.Metadata, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &e, err
}

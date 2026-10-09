package exercise

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, e domain.Exercise) (string, error) {
	query := `
		INSERT INTO exercises (topic_id, type, question, correct_answer, explanation, difficulty, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	metadata := e.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	if !json.Valid([]byte(metadata)) {
		return "", errors.New("invalid metadata JSON")
	}
	err := r.db.QueryRowContext(ctx, query, e.TopicID, e.Type, e.Question, e.CorrectAnswer, e.Explanation, e.Difficulty, metadata).Scan(&id)
	return id, err
}

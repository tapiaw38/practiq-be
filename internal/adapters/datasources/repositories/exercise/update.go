package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, e domain.Exercise) error {
	query := `UPDATE exercises SET type=$1, question=$2, correct_answer=$3, explanation=$4, difficulty=$5, metadata=$6 WHERE id=$7`
	metadata := e.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	_, err := r.db.ExecContext(ctx, query, e.Type, e.Question, e.CorrectAnswer, e.Explanation, e.Difficulty, metadata, id)
	return err
}

package exercise

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, e domain.Exercise) error {
	query := `UPDATE exercises SET type=$1, question=$2, correct_answer=$3, explanation=$4, difficulty=$5, metadata=$6 WHERE id=$7 AND ($8 = '' OR school_id = NULLIF($8, '')::uuid)`
	metadata := e.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	if !json.Valid([]byte(metadata)) {
		return errors.New("invalid metadata JSON")
	}
	_, err := r.db.ExecContext(ctx, query, e.Type, e.Question, e.CorrectAnswer, e.Explanation, e.Difficulty, metadata, id, tenant.SchoolID(ctx))
	return err
}

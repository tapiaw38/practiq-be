package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, s domain.LearningStrategy) error {
	query := `
		UPDATE learning_strategies
		SET name = $1, code = $2, description = $3, status = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, s.Name, s.Code, s.Description, s.Status, id)
	return err
}

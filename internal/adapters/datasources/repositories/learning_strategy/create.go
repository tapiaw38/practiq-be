package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, s domain.LearningStrategy) (string, error) {
	query := `
		INSERT INTO learning_strategies (name, code, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	status := s.Status
	if status == "" {
		status = "active"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, s.Name, s.Code, s.Description, status).Scan(&id)
	return id, err
}

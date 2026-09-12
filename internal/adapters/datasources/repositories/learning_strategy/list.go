package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context) ([]domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE status='active' ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strategies []domain.LearningStrategy
	for rows.Next() {
		var s domain.LearningStrategy
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt); err != nil {
			return nil, err
		}
		strategies = append(strategies, s)
	}
	return strategies, nil
}

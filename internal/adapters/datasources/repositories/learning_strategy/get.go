package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s domain.LearningStrategy
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

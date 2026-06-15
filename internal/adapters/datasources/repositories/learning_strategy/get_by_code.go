package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE code = $1`
	row := r.db.QueryRowContext(ctx, query, code)
	var s domain.LearningStrategy
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

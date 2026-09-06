package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE code = $1 AND ($2 = '' OR school_id = NULLIF($2,'')::uuid)`
	row := r.db.QueryRowContext(ctx, query, code, tenant.SchoolID(ctx))
	var s domain.LearningStrategy
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

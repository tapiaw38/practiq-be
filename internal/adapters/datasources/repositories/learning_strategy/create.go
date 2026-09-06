package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Create(ctx context.Context, s domain.LearningStrategy) (string, error) {
	query := `
		INSERT INTO learning_strategies (name, code, description, status, school_id)
		VALUES ($1, $2, $3, $4, NULLIF($5,'')::uuid)
		RETURNING id
	`
	status := s.Status
	if status == "" {
		status = "active"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, s.Name, s.Code, s.Description, status, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}

package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, s domain.School) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO schools (name, kind, billing, created_by)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id
	`, s.Name, s.Kind, s.Billing, s.CreatedBy).Scan(&id)
	return id, err
}

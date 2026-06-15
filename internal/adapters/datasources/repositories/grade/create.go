package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, grade domain.Grade) (string, error) {
	query := `
		INSERT INTO grades (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, grade.Name, grade.Description, grade.CreatedBy).Scan(&id)
	return id, err
}

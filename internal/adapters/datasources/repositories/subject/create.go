package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, subject domain.Subject) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO subjects (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`, subject.Name, subject.Description, subject.CreatedBy).Scan(&id)
	return id, err
}

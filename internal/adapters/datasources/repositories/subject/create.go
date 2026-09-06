package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Create(ctx context.Context, subject domain.Subject) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO subjects (name, description, created_by, school_id)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid)
		RETURNING id
	`, subject.Name, subject.Description, subject.CreatedBy, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}

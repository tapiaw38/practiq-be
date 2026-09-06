package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, subject domain.Subject) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE subjects SET name = $1, description = $2 WHERE id = $3 AND ($4 = '' OR school_id = NULLIF($4, '')::uuid)
	`, subject.Name, subject.Description, id, tenant.SchoolID(ctx))
	return err
}

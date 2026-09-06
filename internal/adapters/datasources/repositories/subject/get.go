package subject

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
		WHERE id = $1 AND ($2 = '' OR school_id = NULLIF($2, '')::uuid)
	`, id, tenant.SchoolID(ctx)).Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

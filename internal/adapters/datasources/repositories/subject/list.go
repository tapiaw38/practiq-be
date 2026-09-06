package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) List(ctx context.Context) ([]domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
		WHERE ($1 = '' OR school_id = NULLIF($1, '')::uuid)
		ORDER BY name ASC
	`, tenant.SchoolID(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := []domain.Subject{}
	for rows.Next() {
		var subject domain.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

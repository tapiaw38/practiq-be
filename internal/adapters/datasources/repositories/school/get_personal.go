package school

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetPersonal(ctx context.Context, ownerID string) (*domain.School, error) {
	var s domain.School
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, kind, billing, COALESCE(created_by, ''), created_at
		FROM schools
		WHERE created_by = $1 AND kind = 'personal'
		ORDER BY created_at
		LIMIT 1
	`, ownerID).Scan(&s.ID, &s.Name, &s.Kind, &s.Billing, &s.CreatedBy, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		// Not having one is the normal state of a student, not a failure.
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

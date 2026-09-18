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

// CreateWithAdmin keeps an institution and its first administrator atomic.
// An institution with no active admin is unusable by design, not a valid
// intermediate state a partially failed request may leave behind.
func (r *repository) CreateWithAdmin(ctx context.Context, s domain.School, adminID string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var id string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO schools (name, kind, billing, created_by)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id
	`, s.Name, s.Kind, s.Billing, s.CreatedBy).Scan(&id)
	if err != nil {
		return "", err
	}
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO school_members (school_id, user_id, role, active)
		VALUES ($1, $2, 'admin', TRUE)
	`, id, adminID); err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}

package school

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.School, error) {
	var s domain.School
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, kind, billing, status, COALESCE(created_by, ''), created_at,
		       closed_at, COALESCE(closed_by, ''), COALESCE(close_reason, '')
		FROM schools WHERE id = $1
	`, id).Scan(&s.ID, &s.Name, &s.Kind, &s.Billing, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.ClosedAt, &s.ClosedBy, &s.CloseReason)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) List(ctx context.Context) ([]domain.School, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, kind, billing, status, COALESCE(created_by, ''), created_at,
		       closed_at, COALESCE(closed_by, ''), COALESCE(close_reason, '')
		FROM schools ORDER BY kind, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schools := []domain.School{}
	for rows.Next() {
		var s domain.School
		if err := rows.Scan(&s.ID, &s.Name, &s.Kind, &s.Billing, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.ClosedAt, &s.ClosedBy, &s.CloseReason); err != nil {
			return nil, err
		}
		schools = append(schools, s)
	}
	return schools, rows.Err()
}

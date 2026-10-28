package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Close(ctx context.Context, id, closedBy, reason string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET status = 'closed', closed_at = NOW(),
			closed_by = NULLIF($2, ''), close_reason = NULLIF($3, ''), updated_at = NOW()
		WHERE id = $1
	`, id, closedBy, reason)
	return err
}

func (r *repository) Suspend(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET status = 'suspended', updated_at = NOW()
		WHERE id = $1 AND status = 'active'
	`, id)
	return err
}

func (r *repository) Reopen(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET status = 'active', closed_at = NULL, closed_by = NULL,
			close_reason = NULL, updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *repository) Update(ctx context.Context, id string, s domain.School) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET
			name = COALESCE(NULLIF($2, ''), name),
			kind = COALESCE(NULLIF($3, ''), kind),
			billing = COALESCE(NULLIF($4, ''), billing),
			updated_at = NOW()
		WHERE id = $1
	`, id, s.Name, s.Kind, s.Billing)
	return err
}

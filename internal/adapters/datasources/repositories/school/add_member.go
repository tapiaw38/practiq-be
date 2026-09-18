package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// AddMember is an upsert: a school administrator can correct a person's role
// and reactivate a previously suspended membership from the same screen.
func (r *repository) AddMember(ctx context.Context, m domain.SchoolMember) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO school_members (school_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (school_id, user_id) DO UPDATE
		SET role = EXCLUDED.role, active = TRUE
	`, m.SchoolID, m.UserID, m.Role)
	return err
}

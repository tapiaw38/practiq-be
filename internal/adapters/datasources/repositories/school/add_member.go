package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) AddMember(ctx context.Context, m domain.SchoolMember) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO school_members (school_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (school_id, user_id) DO UPDATE
		SET role = EXCLUDED.role, active = TRUE
	`, m.SchoolID, m.UserID, m.Role)
	return err
}

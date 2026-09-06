package studentinvitation

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

// Revoke pide el teacher_id además del id: sin eso, conocer el id de una
// invitación ajena alcanzaría para darla de baja.
func (r *repository) Revoke(ctx context.Context, id, teacherID string) error {
	query := `
		UPDATE student_invitations
		SET revoked_at = NOW()
		WHERE id = $1 AND teacher_id = $2 AND revoked_at IS NULL AND ($3 = '' OR school_id = NULLIF($3, '')::uuid)
	`

	_, err := r.db.ExecContext(ctx, query, id, teacherID, tenant.SchoolID(ctx))

	return err
}

package studentinvitation

import "context"

// Revoke pide el teacher_id además del id: sin eso, conocer el id de una
// invitación ajena alcanzaría para darla de baja.
func (r *repository) Revoke(ctx context.Context, id, teacherID string) error {
	query := `
		UPDATE student_invitations
		SET revoked_at = NOW()
		WHERE id = $1 AND teacher_id = $2 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id, teacherID)

	return err
}

package studentinvitation

import (
	"context"
	"log"
)

// Redeem anota el canje y suma el uso solo si el alumno no había usado antes
// este código. El INSERT hace de candado: si dos pedidos del mismo alumno
// llegan juntos, la clave primaria deja pasar uno solo y el contador no se
// infla.
func (r *repository) Redeem(ctx context.Context, invitationID, studentID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO invitation_redemptions (invitation_id, student_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, invitationID, studentID)
	if err != nil {
		return false, err
	}

	inserted, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if inserted == 0 {
		return false, nil
	}

	if _, err := r.db.ExecContext(ctx, `
		UPDATE student_invitations SET uses = uses + 1 WHERE id = $1
	`, invitationID); err != nil {
		// El contador es informativo y el canje ya quedó anotado. Devolver el
		// error acá haría que el caller aborte antes de crear el vínculo, que
		// es lo único que de verdad importa: se registra y se sigue.
		log.Printf("[invitation] uses counter update failed invitation_id=%s: %v", invitationID, err)
	}

	return true, nil
}

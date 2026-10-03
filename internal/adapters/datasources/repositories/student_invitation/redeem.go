package studentinvitation

import (
	"context"
	"log"
)

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

		log.Printf("[invitation] uses counter update failed invitation_id=%s: %v", invitationID, err)
	}

	return true, nil
}

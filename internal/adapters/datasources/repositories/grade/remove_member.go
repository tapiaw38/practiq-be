package grade

import "context"

func (r *repository) RemoveMember(ctx context.Context, gradeID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM grade_memberships
		WHERE grade_id = $1 AND user_id = $2
	`, gradeID, userID)
	return err
}

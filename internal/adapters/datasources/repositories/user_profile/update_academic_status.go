package userprofile

import "context"

func (r *repository) UpdateAcademicStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE user_profiles
		SET academic_status = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

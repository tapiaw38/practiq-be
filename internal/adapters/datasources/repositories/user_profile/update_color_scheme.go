package userprofile

import "context"

func (r *repository) UpdateColorScheme(ctx context.Context, id, scheme string) error {
	query := `
		UPDATE user_profiles
		SET color_scheme = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, scheme)
	return err
}

package userprofile

import "context"

func (r *repository) UpdateAvatarSeed(ctx context.Context, id, seed string) error {
	query := `
		UPDATE user_profiles
		SET avatar_seed = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, seed)
	return err
}

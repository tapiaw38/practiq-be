package userprofile

import "context"

func (r *repository) UpdateProfileType(ctx context.Context, id, profileType string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE user_profiles SET profile_type = $2 WHERE id = $1`, id, profileType)
	return err
}

package userprofile

import "context"

func (r *repository) UpdateUITheme(ctx context.Context, id, uiTheme string) error {
	query := `
		UPDATE user_profiles
		SET ui_theme = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, uiTheme)
	return err
}

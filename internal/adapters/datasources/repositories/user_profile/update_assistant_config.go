package userprofile

import "context"

func (r *repository) UpdateAssistantConfig(ctx context.Context, id, baseURL, apiKey, uiTheme string) error {
	query := `
		UPDATE user_profiles
		SET assistant_base_url = $2, assistant_api_key = $3, ui_theme = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, baseURL, apiKey, uiTheme)
	return err
}

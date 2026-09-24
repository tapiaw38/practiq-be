package gilliesettings

import (
	"context"
	"database/sql"
)

type Settings struct {
	BaseURL         string
	APIKeyEncrypted string
	APIKeyLast4     string
	UpdatedBy       string
}

type Repository interface {
	Get(context.Context) (Settings, error)
	Save(context.Context, Settings) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }

func (r *repository) Get(ctx context.Context) (v Settings, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT base_url, api_key_encrypted, api_key_last4, updated_by
		FROM gillie_settings WHERE id = 1
	`).Scan(&v.BaseURL, &v.APIKeyEncrypted, &v.APIKeyLast4, &v.UpdatedBy)
	if err == sql.ErrNoRows {
		return Settings{}, nil
	}
	return v, err
}

// Save keeps the stored key when APIKeyEncrypted is empty, so editing the base
// URL alone never needs the secret to travel again.
func (r *repository) Save(ctx context.Context, v Settings) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO gillie_settings (id, base_url, api_key_encrypted, api_key_last4, updated_by, updated_at)
		VALUES (1, $1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET
			base_url          = EXCLUDED.base_url,
			api_key_encrypted = COALESCE(NULLIF(EXCLUDED.api_key_encrypted, ''), gillie_settings.api_key_encrypted),
			api_key_last4     = CASE WHEN EXCLUDED.api_key_encrypted <> ''
			                         THEN EXCLUDED.api_key_last4
			                         ELSE gillie_settings.api_key_last4 END,
			updated_by        = EXCLUDED.updated_by,
			updated_at        = NOW()
	`, v.BaseURL, v.APIKeyEncrypted, v.APIKeyLast4, v.UpdatedBy)
	return err
}

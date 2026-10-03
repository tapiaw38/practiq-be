CREATE TABLE IF NOT EXISTS gillie_settings (
    id                SMALLINT PRIMARY KEY CHECK (id = 1),
    base_url          TEXT NOT NULL DEFAULT '',
    api_key_encrypted TEXT NOT NULL DEFAULT '',
    api_key_last4     TEXT NOT NULL DEFAULT '',
    updated_by        TEXT NOT NULL DEFAULT '',
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO gillie_settings (id) VALUES (1) ON CONFLICT (id) DO NOTHING;

UPDATE user_profiles
SET assistant_api_key = '', assistant_base_url = ''
WHERE assistant_api_key <> '' OR assistant_base_url <> '';

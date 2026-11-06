ALTER TABLE user_profiles
    ADD COLUMN IF NOT EXISTS color_scheme VARCHAR(8) NOT NULL DEFAULT 'light';

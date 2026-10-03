ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS ui_theme VARCHAR(30) NOT NULL DEFAULT 'primary' CHECK (ui_theme IN ('primary', 'secondary'));

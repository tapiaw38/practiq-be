ALTER TABLE grades
    ADD COLUMN IF NOT EXISTS visual_theme VARCHAR(30) NOT NULL DEFAULT 'primary'
        CHECK (visual_theme IN ('primary', 'secondary'));

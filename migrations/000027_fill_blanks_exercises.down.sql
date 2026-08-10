-- Refuse to roll back while fill_blanks exercises exist: dropping them would
-- destroy content a teacher wrote. Convert or delete them deliberately first.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM exercises WHERE type = 'fill_blanks') THEN
        RAISE EXCEPTION
            'there are fill_blanks exercises; convert or remove them before rolling back this migration';
    END IF;
END $$;

ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_type_check;
ALTER TABLE exercises ADD CONSTRAINT exercises_type_check
    CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas', 'attachment'));

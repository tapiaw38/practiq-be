ALTER TABLE practice_sheets
    ADD COLUMN IF NOT EXISTS sheet_type VARCHAR(30) NOT NULL DEFAULT 'practice'
        CHECK (sheet_type IN ('practice', 'level_test'));
ALTER TABLE practice_sheets ADD COLUMN IF NOT EXISTS test_style VARCHAR(20) NOT NULL DEFAULT 'keyboard' CHECK (test_style IN ('keyboard', 'canvas'));

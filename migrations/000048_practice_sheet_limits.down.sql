UPDATE level_test_submissions SET submitted_at = NOW() WHERE submitted_at IS NULL;
ALTER TABLE level_test_submissions ALTER COLUMN submitted_at SET NOT NULL;

ALTER TABLE level_test_submissions
    DROP COLUMN IF EXISTS started_at,
    DROP COLUMN IF EXISTS attempts;

ALTER TABLE practice_sheets
    DROP CONSTRAINT IF EXISTS practice_sheets_time_limit_positive,
    DROP CONSTRAINT IF EXISTS practice_sheets_max_attempts_positive;

ALTER TABLE practice_sheets
    DROP COLUMN IF EXISTS time_limit_minutes,
    DROP COLUMN IF EXISTS max_attempts;

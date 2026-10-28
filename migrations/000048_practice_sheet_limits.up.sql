-- NULL means "no limit set". A level test still allows one attempt when
-- max_attempts is NULL, which is what it allowed before this column existed.
ALTER TABLE practice_sheets
    ADD COLUMN IF NOT EXISTS max_attempts INTEGER,
    ADD COLUMN IF NOT EXISTS time_limit_minutes INTEGER;

ALTER TABLE practice_sheets
    ADD CONSTRAINT practice_sheets_max_attempts_positive
        CHECK (max_attempts IS NULL OR max_attempts > 0),
    ADD CONSTRAINT practice_sheets_time_limit_positive
        CHECK (time_limit_minutes IS NULL OR time_limit_minutes > 0);

-- attempts counts how many times the student has submitted. Existing rows are
-- one submission each, which is what the table could hold until now.
-- started_at is when the student first opened the test; the time limit runs
-- from there, and it stays NULL for everyone who opened one before this.
ALTER TABLE level_test_submissions
    ADD COLUMN IF NOT EXISTS attempts INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ;

-- A row now exists from the moment a student opens the test, before anything
-- has been submitted, so the column stops being a lie for that window.
ALTER TABLE level_test_submissions ALTER COLUMN submitted_at DROP NOT NULL;

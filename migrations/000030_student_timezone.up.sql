-- The streak compares calendar days, and both sides of that comparison were
-- wrong for a UTC-3 audience: the UTC day flips at 21:00 local, so two
-- sessions the same evening counted as two days, and two consecutive evenings
-- counted as one.
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS timezone TEXT NOT NULL DEFAULT '';

-- TIMESTAMP without a zone stores an ambiguous instant: the value's meaning
-- depends on the session timezone that happened to be active when it was
-- written. TIMESTAMPTZ pins the instant and converts on read.
--
-- Existing rows were written by a UTC container, so they are reinterpreted as
-- UTC, which is what they already were.
ALTER TABLE student_topic_progress
    ALTER COLUMN last_practiced_at TYPE TIMESTAMPTZ USING last_practiced_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';

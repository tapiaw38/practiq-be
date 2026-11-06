ALTER TABLE student_topic_progress
    ALTER COLUMN last_practiced_at TYPE TIMESTAMP USING last_practiced_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'UTC';

ALTER TABLE user_profiles DROP COLUMN IF EXISTS timezone;

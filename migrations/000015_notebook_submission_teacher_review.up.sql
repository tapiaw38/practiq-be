ALTER TABLE notebook_submissions
ADD COLUMN IF NOT EXISTS teacher_is_correct BOOLEAN,
ADD COLUMN IF NOT EXISTS teacher_feedback TEXT NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS teacher_reviewed_at TIMESTAMP;

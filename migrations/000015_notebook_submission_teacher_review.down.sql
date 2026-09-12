ALTER TABLE notebook_submissions
DROP COLUMN IF EXISTS teacher_reviewed_at,
DROP COLUMN IF EXISTS teacher_feedback,
DROP COLUMN IF EXISTS teacher_is_correct;

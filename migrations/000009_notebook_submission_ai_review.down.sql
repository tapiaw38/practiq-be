ALTER TABLE notebook_submissions
DROP COLUMN IF EXISTS ai_reviewed_at,
DROP COLUMN IF EXISTS ai_feedback,
DROP COLUMN IF EXISTS ai_is_correct;

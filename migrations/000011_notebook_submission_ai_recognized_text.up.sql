ALTER TABLE notebook_submissions
ADD COLUMN IF NOT EXISTS ai_recognized_text TEXT NOT NULL DEFAULT '';

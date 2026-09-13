-- The assistant's verdict on an attachment is a suggestion for the teacher, not
-- the grade: is_correct stays false until a human confirms it.
ALTER TABLE student_attempts ADD COLUMN IF NOT EXISTS ai_is_correct BOOLEAN;

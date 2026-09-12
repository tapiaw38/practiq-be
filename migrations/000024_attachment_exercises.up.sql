-- Exercises answered by uploading a file (audio, pdf, image, document).
-- The accepted formats live in the exercise metadata, e.g. {"accept":["audio","pdf"]}.
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_type_check;
ALTER TABLE exercises ADD CONSTRAINT exercises_type_check
    CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas', 'attachment'));

ALTER TABLE student_attempts
    ADD COLUMN IF NOT EXISTS attachment_url TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS attachment_name TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS attachment_content_type VARCHAR(150) DEFAULT '',
    -- Set when the assistant could not evaluate the file, so the teacher must.
    ADD COLUMN IF NOT EXISTS needs_teacher_review BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS teacher_is_correct BOOLEAN,
    ADD COLUMN IF NOT EXISTS teacher_feedback TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS teacher_reviewed_at TIMESTAMP;

-- Drives the teacher's pending-review queue.
CREATE INDEX IF NOT EXISTS idx_student_attempts_pending_review
    ON student_attempts(created_at DESC)
    WHERE needs_teacher_review AND teacher_reviewed_at IS NULL;

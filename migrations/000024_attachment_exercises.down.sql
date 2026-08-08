DROP INDEX IF EXISTS idx_student_attempts_pending_review;

ALTER TABLE student_attempts
    DROP COLUMN IF EXISTS attachment_url,
    DROP COLUMN IF EXISTS attachment_name,
    DROP COLUMN IF EXISTS attachment_content_type,
    DROP COLUMN IF EXISTS needs_teacher_review,
    DROP COLUMN IF EXISTS teacher_is_correct,
    DROP COLUMN IF EXISTS teacher_feedback,
    DROP COLUMN IF EXISTS teacher_reviewed_at;

-- Attempts referencing attachment exercises must go before the constraint is
-- narrowed again, otherwise the exercises would violate it.
DELETE FROM exercises WHERE type = 'attachment';
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_type_check;
ALTER TABLE exercises ADD CONSTRAINT exercises_type_check
    CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas'));

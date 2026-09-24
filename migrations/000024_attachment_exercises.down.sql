-- Refuse to roll back while attachment exercises exist, the same way 000027
-- does for fill_blanks. The previous version ran `DELETE FROM exercises WHERE
-- type = 'attachment'`, and the ON DELETE CASCADE foreign keys took the
-- practice-sheet associations and every student attempt with them: an
-- operational rollback silently destroyed teacher content and student results.
-- Convert or remove them deliberately first.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM exercises WHERE type = 'attachment'
    ) OR EXISTS (
        SELECT 1
        FROM student_attempts
        WHERE attachment_url <> ''
           OR attachment_name <> ''
           OR attachment_content_type <> ''
           OR needs_teacher_review
           OR teacher_is_correct IS NOT NULL
           OR teacher_feedback <> ''
           OR teacher_reviewed_at IS NOT NULL
    ) THEN
        RAISE EXCEPTION
            'attachment exercises or attachment review data exist; remove them before rolling back this migration';
    END IF;
END $$;

DROP INDEX IF EXISTS idx_student_attempts_pending_review;

ALTER TABLE student_attempts
    DROP COLUMN IF EXISTS attachment_url,
    DROP COLUMN IF EXISTS attachment_name,
    DROP COLUMN IF EXISTS attachment_content_type,
    DROP COLUMN IF EXISTS needs_teacher_review,
    DROP COLUMN IF EXISTS teacher_is_correct,
    DROP COLUMN IF EXISTS teacher_feedback,
    DROP COLUMN IF EXISTS teacher_reviewed_at;

ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_type_check;
ALTER TABLE exercises ADD CONSTRAINT exercises_type_check
    CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas'));

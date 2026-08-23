-- Only level tests and the homework notebook are corrected by the teacher.
-- Practice answers were being queued for review, which left the student waiting
-- on a correction that a practice is not supposed to need. Clear the flag on
-- the practice answers already stored so they stop showing up as pending.
UPDATE student_attempts sa
SET needs_teacher_review = FALSE
FROM practice_sheets ps
WHERE ps.id = sa.practice_sheet_id
  AND COALESCE(ps.sheet_type, '') <> 'level_test'
  AND sa.needs_teacher_review;

-- Nobody put a verdict on this answer: the assistant could not read it, or the
-- exercise carried image/audio in its statement. It is left out of the score
-- instead of counted as wrong, and until now that decision lived only in the
-- submit response — every later reader saw is_correct=false, score=0 and
-- reported it as an error the student made.
--
-- Distinct from needs_teacher_review, which says a person still has to settle
-- it. Only a level test sets that; a practice is never queued, so on a practice
-- this is the only trace that the answer was skipped rather than failed.
ALTER TABLE student_attempts
    ADD COLUMN IF NOT EXISTS not_graded BOOLEAN NOT NULL DEFAULT FALSE;

-- Everything still waiting for a teacher is by definition ungraded.
UPDATE student_attempts
SET not_graded = TRUE
WHERE needs_teacher_review AND teacher_reviewed_at IS NULL;

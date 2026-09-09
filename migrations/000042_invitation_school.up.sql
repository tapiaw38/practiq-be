-- An invitation now says which school the student joins.
--
-- It was deduced from the teacher's personal school, which is right for an
-- independent teacher and wrong for one who works at an institution: their
-- students landed in the teacher's own school instead of the institution's.

ALTER TABLE student_invitations
    ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id) ON DELETE CASCADE;

-- Existing invitations belong to the school their teacher administers, which
-- until now was the only school anybody had.
UPDATE student_invitations si
SET school_id = s.id
FROM schools s
WHERE s.created_by = si.teacher_id AND si.school_id IS NULL;

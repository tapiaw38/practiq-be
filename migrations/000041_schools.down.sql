-- The duplicated grades and subjects are not undone: reversing the split would
-- have to guess which copy the original was, and courses now point at the
-- copies. Going back means restoring the global uniqueness by hand.

DROP INDEX IF EXISTS idx_subjects_school_name;
DROP INDEX IF EXISTS idx_grades_school_name;

ALTER TABLE subjects DROP COLUMN IF EXISTS school_id;
ALTER TABLE grades DROP COLUMN IF EXISTS school_id;

DROP INDEX IF EXISTS idx_school_members_user;
DROP TABLE IF EXISTS school_members;
DROP TABLE IF EXISTS schools;

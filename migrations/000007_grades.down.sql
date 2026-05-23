DROP INDEX IF EXISTS idx_courses_grade_id;
DROP INDEX IF EXISTS idx_grade_memberships_user_id;

ALTER TABLE courses
    DROP COLUMN IF EXISTS grade_id;

DROP TABLE IF EXISTS grade_memberships;
DROP TABLE IF EXISTS grades;

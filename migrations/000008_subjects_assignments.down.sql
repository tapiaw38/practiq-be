DROP INDEX IF EXISTS idx_courses_subject_id;
DROP INDEX IF EXISTS idx_assignments_student;
DROP INDEX IF EXISTS idx_assignments_teacher;

ALTER TABLE courses
    DROP COLUMN IF EXISTS subject_id;

DROP TABLE IF EXISTS teacher_student_assignments;

ALTER TABLE user_profiles
    DROP COLUMN IF EXISTS academic_status;

DROP TABLE IF EXISTS subjects;

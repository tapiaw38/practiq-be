-- Remove student_id column from submit_jobs table
DROP INDEX IF EXISTS idx_submit_jobs_student_id;
ALTER TABLE submit_jobs DROP COLUMN IF EXISTS student_id;

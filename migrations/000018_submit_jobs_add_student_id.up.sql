-- Add student_id column to submit_jobs table for ownership verification
ALTER TABLE submit_jobs ADD COLUMN student_id TEXT;

-- Add index for efficient querying by student_id
CREATE INDEX idx_submit_jobs_student_id ON submit_jobs(student_id);

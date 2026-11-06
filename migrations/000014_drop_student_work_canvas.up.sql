-- Add image_url column to student_attempts if not exists
ALTER TABLE student_attempts ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT '';

-- Backfill image_url from student_work_canvas to student_attempts before dropping
UPDATE student_attempts sa
SET image_url = swc.image_url
FROM student_work_canvas swc
WHERE sa.id = swc.attempt_id AND swc.image_url IS NOT NULL AND swc.image_url != '';

-- Drop student_work_canvas table as it is no longer used.
-- Canvas data is stored directly in student_attempts.answer_text field.
DROP TABLE IF EXISTS student_work_canvas;

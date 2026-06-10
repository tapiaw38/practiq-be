-- Drop student_work_canvas table as it is no longer used.
-- Canvas data is stored directly in student_attempts.answer_text field.
DROP TABLE IF EXISTS student_work_canvas;

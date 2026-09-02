ALTER TABLE notebook_pages
DROP COLUMN IF EXISTS statement_text,
DROP COLUMN IF EXISTS statement_verified;

ALTER TABLE notebook_submissions
DROP COLUMN IF EXISTS needs_teacher_review;

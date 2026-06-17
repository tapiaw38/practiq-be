-- Drop indexes
DROP INDEX IF EXISTS idx_courses_deleted_at;
DROP INDEX IF EXISTS idx_notebooks_deleted_at;

-- Remove deleted_at columns
ALTER TABLE courses DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE notebooks DROP COLUMN IF EXISTS deleted_at;

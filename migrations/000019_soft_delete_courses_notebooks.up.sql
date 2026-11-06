-- Add deleted_at column to courses table for soft delete
ALTER TABLE courses ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Add deleted_at column to notebooks table for soft delete
ALTER TABLE notebooks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Create index for filtering out deleted records
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON courses(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_notebooks_deleted_at ON notebooks(deleted_at) WHERE deleted_at IS NULL;

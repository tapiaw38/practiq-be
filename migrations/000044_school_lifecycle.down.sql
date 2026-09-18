DROP INDEX IF EXISTS idx_schools_status;

ALTER TABLE schools
    DROP COLUMN IF EXISTS close_reason,
    DROP COLUMN IF EXISTS closed_by,
    DROP COLUMN IF EXISTS closed_at,
    DROP COLUMN IF EXISTS status;

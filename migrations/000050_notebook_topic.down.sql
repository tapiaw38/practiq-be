DROP INDEX IF EXISTS idx_notebooks_topic_id;
ALTER TABLE notebooks DROP COLUMN IF EXISTS topic_id;

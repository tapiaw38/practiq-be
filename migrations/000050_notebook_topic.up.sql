ALTER TABLE notebooks
    ADD COLUMN IF NOT EXISTS topic_id UUID REFERENCES topics(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_notebooks_topic_id ON notebooks(topic_id);

-- Change unique constraint from (student_id, topic_id, strategy_id) to (student_id, topic_id)
-- This allows proper upsert by student+topic regardless of strategy

-- Drop existing constraint
ALTER TABLE student_topic_progress DROP CONSTRAINT IF EXISTS student_topic_progress_student_id_topic_id_strategy_id_key;

-- Deduplicate: the old constraint included strategy_id (nullable), so rows with
-- NULL strategy_id never conflicted and duplicates accumulated. Keep the most
-- recent row per (student_id, topic_id).
DELETE FROM student_topic_progress a
USING student_topic_progress b
WHERE a.student_id = b.student_id
  AND a.topic_id = b.topic_id
  AND (a.updated_at < b.updated_at
       OR (a.updated_at = b.updated_at AND a.id < b.id));

-- Add new unique constraint on (student_id, topic_id) only
ALTER TABLE student_topic_progress ADD CONSTRAINT student_topic_progress_student_id_topic_id_key UNIQUE (student_id, topic_id);

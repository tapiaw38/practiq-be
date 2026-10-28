-- Revert unique constraint change
ALTER TABLE student_topic_progress DROP CONSTRAINT IF EXISTS student_topic_progress_student_id_topic_id_key;

-- Restore original constraint
ALTER TABLE student_topic_progress ADD CONSTRAINT student_topic_progress_student_id_topic_id_strategy_id_key UNIQUE (student_id, topic_id, strategy_id);

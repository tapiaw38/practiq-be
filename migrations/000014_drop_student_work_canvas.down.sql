-- Recreate student_work_canvas table (reverts 000014_drop_student_work_canvas.up.sql)
CREATE TABLE IF NOT EXISTS student_work_canvas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES student_attempts(id) ON DELETE CASCADE,
    canvas_data JSONB DEFAULT '{}',
    image_url TEXT,
    recognized_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- The up migration moved every canvas image into student_attempts.image_url
-- and then dropped this table. Recreating it empty made the rollback silently
-- destroy every canvas submission, and rolling back past 000017 (which drops
-- image_url) made that loss permanent. Copy the data back.
--
-- canvas_data and recognized_text are not restored: the up migration never
-- carried them across, so they are already gone by this point.
INSERT INTO student_work_canvas (attempt_id, image_url)
SELECT sa.id, sa.image_url
FROM student_attempts sa
WHERE sa.image_url IS NOT NULL AND sa.image_url != '';

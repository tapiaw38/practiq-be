-- Recreate student_work_canvas table (reverts 000014_drop_student_work_canvas.up.sql)
CREATE TABLE IF NOT EXISTS student_work_canvas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES student_attempts(id) ON DELETE CASCADE,
    canvas_data JSONB DEFAULT '{}',
    image_url TEXT,
    recognized_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

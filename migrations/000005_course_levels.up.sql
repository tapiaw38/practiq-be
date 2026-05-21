-- Add level to notebooks
ALTER TABLE notebooks ADD COLUMN IF NOT EXISTS level INT NOT NULL DEFAULT 1;

-- Track student's current level per course
CREATE TABLE IF NOT EXISTS student_course_progress (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id  VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    current_level INT NOT NULL DEFAULT 1,
    updated_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_scp_student_course ON student_course_progress(student_id, course_id);

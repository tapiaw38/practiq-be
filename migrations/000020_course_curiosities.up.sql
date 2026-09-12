CREATE TABLE IF NOT EXISTS course_curiosities (
    course_id UUID PRIMARY KEY REFERENCES courses(id) ON DELETE CASCADE,
    curiosities JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_course_curiosities_course_id ON course_curiosities(course_id);

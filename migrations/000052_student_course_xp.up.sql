-- XP belongs to a student's course journey. The balance is cheap to read;
-- events make every award auditable and, importantly, idempotent.
CREATE TABLE student_course_xp (
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    total_xp INTEGER NOT NULL DEFAULT 0 CHECK (total_xp >= 0),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (student_id, course_id)
);

CREATE TABLE student_course_xp_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    event_key VARCHAR(255) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    points INTEGER NOT NULL CHECK (points > 0),
    practice_sheet_id UUID REFERENCES practice_sheets(id) ON DELETE SET NULL,
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, course_id, event_key)
);

CREATE INDEX idx_student_course_xp_events_course_student
    ON student_course_xp_events(course_id, student_id, created_at DESC);

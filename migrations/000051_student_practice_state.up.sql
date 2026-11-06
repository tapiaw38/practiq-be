-- "Continue" is learning state, not a browser preference. Keep it on the
-- server so it follows the student to every device.
CREATE TABLE student_practice_state (
    student_id VARCHAR(255) PRIMARY KEY REFERENCES user_profiles(id) ON DELETE CASCADE,
    practice_sheet_id UUID REFERENCES practice_sheets(id) ON DELETE SET NULL,
    last_opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_student_practice_state_last_opened
    ON student_practice_state(last_opened_at DESC);

-- Preserve a meaningful Continue action for existing students, without ever
-- trusting a former browser-local value. Level tests use their own flow.
INSERT INTO student_practice_state (student_id, practice_sheet_id, last_opened_at, updated_at)
SELECT DISTINCT ON (sa.student_id)
    sa.student_id,
    sa.practice_sheet_id,
    sa.created_at,
    sa.created_at
FROM student_attempts sa
JOIN practice_sheets ps ON ps.id = sa.practice_sheet_id
JOIN courses c ON c.id = ps.course_id AND c.deleted_at IS NULL
WHERE ps.sheet_type = 'practice'
ORDER BY sa.student_id, sa.created_at DESC
ON CONFLICT (student_id) DO NOTHING;

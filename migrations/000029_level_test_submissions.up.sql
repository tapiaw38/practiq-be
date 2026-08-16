CREATE TABLE IF NOT EXISTS level_test_submissions (
    practice_sheet_id UUID NOT NULL REFERENCES practice_sheets(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (practice_sheet_id, student_id)
);

-- Existing submissions are final too; without this backfill they could be
-- resent once immediately after deployment.
INSERT INTO level_test_submissions (practice_sheet_id, student_id, submitted_at)
SELECT sa.practice_sheet_id, sa.student_id, MIN(sa.created_at)
FROM student_attempts sa
JOIN practice_sheets ps ON ps.id = sa.practice_sheet_id
WHERE ps.sheet_type = 'level_test'
  AND sa.practice_sheet_id IS NOT NULL
GROUP BY sa.practice_sheet_id, sa.student_id
ON CONFLICT (practice_sheet_id, student_id) DO NOTHING;

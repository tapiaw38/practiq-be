-- Course levels expose one promotion test. Multiple rows made the selected
-- test depend on query order and allowed a second test to promote twice.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM practice_sheets
        WHERE sheet_type = 'level_test'
        GROUP BY course_id, level
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'duplicate level tests exist; resolve them before migration 000047';
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS ux_practice_sheets_one_level_test
    ON practice_sheets (course_id, level)
    WHERE sheet_type = 'level_test';

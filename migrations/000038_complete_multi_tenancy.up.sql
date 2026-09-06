-- Complete tenant ownership for records not covered by 000037. Existing data
-- is assigned to one explicit legacy school before tenant keys become required.
INSERT INTO schools (name, slug, created_by)
SELECT 'Legacy school', 'legacy', 'migration'
WHERE EXISTS (SELECT 1 FROM user_profiles)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO school_memberships (school_id, user_id, membership_role, profile_type)
SELECT s.id, up.id,
       CASE WHEN up.profile_type = 'teacher' THEN 'admin' ELSE 'member' END,
       up.profile_type
FROM user_profiles up
JOIN schools s ON s.slug = 'legacy'
ON CONFLICT (school_id, user_id) DO NOTHING;

DO $$
DECLARE legacy_school_id UUID;
BEGIN
    SELECT id INTO legacy_school_id FROM schools WHERE slug = 'legacy';
    IF legacy_school_id IS NULL THEN RETURN; END IF;

    ALTER TABLE teacher_student_assignments ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id);
    ALTER TABLE course_learning_strategies ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id);
    ALTER TABLE notebook_pages ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id);
    ALTER TABLE notebook_submissions ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id);
    ALTER TABLE level_test_submissions ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id);

    UPDATE teacher_student_assignments SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE course_learning_strategies cls SET school_id = c.school_id FROM courses c WHERE cls.course_id = c.id AND cls.school_id IS NULL;
    UPDATE notebook_pages np SET school_id = n.school_id FROM notebooks n WHERE np.notebook_id = n.id AND np.school_id IS NULL;
    UPDATE notebook_submissions ns SET school_id = np.school_id FROM notebook_pages np WHERE ns.page_id = np.id AND ns.school_id IS NULL;
    UPDATE level_test_submissions lts SET school_id = ps.school_id FROM practice_sheets ps WHERE lts.practice_sheet_id = ps.id AND lts.school_id IS NULL;

    -- Root records lacking a parent belong to legacy school; child records get
    -- their parent's tenant so cross-school IDs cannot be reintroduced.
    UPDATE grades SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE subjects SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE courses SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE materials m SET school_id = c.school_id FROM courses c WHERE m.course_id = c.id AND m.school_id IS NULL;
    UPDATE topics t SET school_id = c.school_id FROM courses c WHERE t.course_id = c.id AND t.school_id IS NULL;
    UPDATE exercises e SET school_id = t.school_id FROM topics t WHERE e.topic_id = t.id AND e.school_id IS NULL;
    UPDATE practice_sheets ps SET school_id = c.school_id FROM courses c WHERE ps.course_id = c.id AND ps.school_id IS NULL;
    UPDATE notebooks n SET school_id = c.school_id FROM courses c WHERE n.course_id = c.id AND n.school_id IS NULL;
    UPDATE enrollments e SET school_id = c.school_id FROM courses c WHERE e.course_id = c.id AND e.school_id IS NULL;
    UPDATE student_attempts sa SET school_id = e.school_id FROM exercises e WHERE sa.exercise_id = e.id AND sa.school_id IS NULL;
    UPDATE student_topic_progress stp SET school_id = t.school_id FROM topics t WHERE stp.topic_id = t.id AND stp.school_id IS NULL;
    UPDATE student_course_progress scp SET school_id = c.school_id FROM courses c WHERE scp.course_id = c.id AND scp.school_id IS NULL;
    UPDATE ai_conversations SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE ai_help_requests SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE student_invitations SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE notifications SET school_id = legacy_school_id WHERE school_id IS NULL;
    UPDATE learning_strategies SET school_id = legacy_school_id WHERE school_id IS NULL;

    ALTER TABLE teacher_student_assignments ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE course_learning_strategies ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE notebook_pages ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE notebook_submissions ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE level_test_submissions ALTER COLUMN school_id SET NOT NULL;

    ALTER TABLE grades ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE subjects ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE courses ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE materials ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE topics ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE exercises ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE practice_sheets ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE notebooks ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE enrollments ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE student_attempts ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE student_topic_progress ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE student_course_progress ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE ai_conversations ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE ai_help_requests ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE student_invitations ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE notifications ALTER COLUMN school_id SET NOT NULL;
    ALTER TABLE learning_strategies ALTER COLUMN school_id SET NOT NULL;
END $$;

CREATE INDEX IF NOT EXISTS idx_teacher_student_assignments_school ON teacher_student_assignments(school_id);
CREATE INDEX IF NOT EXISTS idx_course_learning_strategies_school ON course_learning_strategies(school_id);
CREATE INDEX IF NOT EXISTS idx_notebook_pages_school ON notebook_pages(school_id);
CREATE INDEX IF NOT EXISTS idx_notebook_submissions_school ON notebook_submissions(school_id);
CREATE INDEX IF NOT EXISTS idx_level_test_submissions_school ON level_test_submissions(school_id);

-- Existing pair uniqueness is global. Tenant ownership makes same teacher and
-- student valid in separate schools, so replace it with a tenant-aware key.
ALTER TABLE teacher_student_assignments DROP CONSTRAINT IF EXISTS teacher_student_assignments_teacher_id_student_id_key;
CREATE UNIQUE INDEX IF NOT EXISTS teacher_student_assignments_school_teacher_student_key
    ON teacher_student_assignments(school_id, teacher_id, student_id);

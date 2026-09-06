CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS school_memberships (
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    membership_role VARCHAR(20) NOT NULL DEFAULT 'member'
        CHECK (membership_role IN ('admin', 'member')),
    profile_type VARCHAR(20) NOT NULL DEFAULT 'student'
        CHECK (profile_type IN ('teacher', 'student')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (school_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_school_memberships_user
    ON school_memberships(user_id);

ALTER TABLE user_profiles
    ADD COLUMN IF NOT EXISTS profile_type VARCHAR(30);

UPDATE user_profiles
SET profile_type = 'student'
WHERE profile_type IS NULL OR profile_type NOT IN ('teacher', 'student');

ALTER TABLE user_profiles
    ALTER COLUMN profile_type SET DEFAULT 'student';

ALTER TABLE user_profiles
    ALTER COLUMN profile_type SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_school_memberships_school_role
    ON school_memberships(school_id, membership_role);

-- Tenant key is intentionally nullable during rollout. Existing test data can
-- remain available while new writes are migrated to an active school.
DO $$
DECLARE
    table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'user_profiles', 'grades', 'subjects', 'courses', 'materials',
        'topics', 'exercises', 'practice_sheets', 'notebooks',
        'enrollments', 'student_attempts', 'student_topic_progress',
        'student_course_progress', 'ai_conversations', 'ai_help_requests',
        'student_invitations', 'notifications', 'learning_strategies'
    ] LOOP
        EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id)', table_name);
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_school_id ON %I(school_id)', table_name, table_name);
    END LOOP;
END $$;

-- Child records inherit tenant from their course graph. This keeps writes from
-- creating unscoped records when legacy handlers do not yet send school_id.
CREATE OR REPLACE FUNCTION set_practiq_school_id_from_parent()
RETURNS trigger AS $$
BEGIN
    IF NEW.school_id IS NULL THEN
        IF TG_TABLE_NAME IN ('materials', 'topics', 'notebooks', 'practice_sheets', 'enrollments') THEN
            EXECUTE format('SELECT school_id FROM courses WHERE id = $1') INTO NEW.school_id USING NEW.course_id;
        ELSIF TG_TABLE_NAME = 'exercises' THEN
            SELECT c.school_id INTO NEW.school_id
            FROM topics t JOIN courses c ON c.id = t.course_id
            WHERE t.id = NEW.topic_id;
        ELSIF TG_TABLE_NAME IN ('student_course_progress') THEN
            SELECT school_id INTO NEW.school_id FROM courses WHERE id = NEW.course_id;
        ELSIF TG_TABLE_NAME = 'student_topic_progress' THEN
            SELECT c.school_id INTO NEW.school_id
            FROM topics t JOIN courses c ON c.id = t.course_id WHERE t.id = NEW.topic_id;
        ELSIF TG_TABLE_NAME = 'student_attempts' THEN
            SELECT c.school_id INTO NEW.school_id
            FROM exercises e JOIN topics t ON t.id = e.topic_id JOIN courses c ON c.id = t.course_id
            WHERE e.id = NEW.exercise_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_materials_school_id ON materials;
CREATE TRIGGER trg_materials_school_id BEFORE INSERT ON materials FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_topics_school_id ON topics;
CREATE TRIGGER trg_topics_school_id BEFORE INSERT ON topics FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_exercises_school_id ON exercises;
CREATE TRIGGER trg_exercises_school_id BEFORE INSERT ON exercises FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_notebooks_school_id ON notebooks;
CREATE TRIGGER trg_notebooks_school_id BEFORE INSERT ON notebooks FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_practice_sheets_school_id ON practice_sheets;
CREATE TRIGGER trg_practice_sheets_school_id BEFORE INSERT ON practice_sheets FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_enrollments_school_id ON enrollments;
CREATE TRIGGER trg_enrollments_school_id BEFORE INSERT ON enrollments FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_student_course_progress_school_id ON student_course_progress;
CREATE TRIGGER trg_student_course_progress_school_id BEFORE INSERT ON student_course_progress FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_student_topic_progress_school_id ON student_topic_progress;
CREATE TRIGGER trg_student_topic_progress_school_id BEFORE INSERT ON student_topic_progress FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();
DROP TRIGGER IF EXISTS trg_student_attempts_school_id ON student_attempts;
CREATE TRIGGER trg_student_attempts_school_id BEFORE INSERT ON student_attempts FOR EACH ROW EXECUTE FUNCTION set_practiq_school_id_from_parent();

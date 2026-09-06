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
        EXECUTE format('ALTER TABLE %I DROP COLUMN IF EXISTS school_id', table_name);
    END LOOP;
END $$;

DROP TRIGGER IF EXISTS trg_materials_school_id ON materials;
DROP TRIGGER IF EXISTS trg_topics_school_id ON topics;
DROP TRIGGER IF EXISTS trg_exercises_school_id ON exercises;
DROP TRIGGER IF EXISTS trg_notebooks_school_id ON notebooks;
DROP TRIGGER IF EXISTS trg_practice_sheets_school_id ON practice_sheets;
DROP TRIGGER IF EXISTS trg_enrollments_school_id ON enrollments;
DROP TRIGGER IF EXISTS trg_student_course_progress_school_id ON student_course_progress;
DROP TRIGGER IF EXISTS trg_student_topic_progress_school_id ON student_topic_progress;
DROP TRIGGER IF EXISTS trg_student_attempts_school_id ON student_attempts;
DROP FUNCTION IF EXISTS set_practiq_school_id_from_parent();

DROP TABLE IF EXISTS school_memberships;
DROP TABLE IF EXISTS schools;

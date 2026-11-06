-- Schools: the academic tree and the people in it stop being one shared space.
--
-- This migration only moves data into shape. Nothing filters by school yet, so
-- behaviour is unchanged and the result can be verified before anything starts
-- depending on it.

CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    -- personal: one teacher who signed up on their own and administers it.
    -- institution: created by hand, has admins and teachers.
    kind VARCHAR(20) NOT NULL DEFAULT 'personal' CHECK (kind IN ('personal', 'institution')),
    -- direct means invoiced outside the product: no entitlement is consulted
    -- and no student limit applies.
    billing VARCHAR(20) NOT NULL DEFAULT 'subscription' CHECK (billing IN ('subscription', 'direct')),
    -- Who it was created for. Personal schools point at their teacher; an
    -- institution created by an operator has no single owner.
    created_by VARCHAR(255) REFERENCES user_profiles(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Membership carries the role, so one person can be an admin of their own
-- school and a teacher at an institution without either row knowing about the
-- other. This is also what makes "is X an admin of Y" have a single answer.
CREATE TABLE IF NOT EXISTS school_members (
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
    -- Deactivated members keep their history and lose access to this school's
    -- courses. Used when a downgrade leaves more students than the plan allows.
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (school_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_school_members_user ON school_members(user_id);

ALTER TABLE grades ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id) ON DELETE CASCADE;
ALTER TABLE subjects ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id) ON DELETE CASCADE;

-- The old constraint is the shared catalogue written into the schema: two
-- schools cannot both have a "5to Grado" while the name is globally unique.
ALTER TABLE grades DROP CONSTRAINT IF EXISTS grades_name_key;
ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_name_key;

-- ── Backfill ────────────────────────────────────────────────────────────────

-- One personal school per teacher. The name is a placeholder: teacher names
-- live in auth-api-be since migration 000036, so this migration cannot know
-- them. The application sets a real one from the teacher's name on signup, and
-- existing teachers can rename theirs.
INSERT INTO schools (name, kind, billing, created_by)
SELECT 'Mi escuela', 'personal', 'subscription', up.id
FROM user_profiles up
WHERE up.profile_type = 'teacher'
  AND NOT EXISTS (SELECT 1 FROM schools s WHERE s.created_by = up.id);

INSERT INTO school_members (school_id, user_id, role)
SELECT s.id, s.created_by, 'admin'
FROM schools s
WHERE s.created_by IS NOT NULL
ON CONFLICT DO NOTHING;

-- Grades and subjects move into the school of the teacher whose courses use
-- them. One is claimed outright; a second school needing the same one gets its
-- own copy, because a row cannot belong to two schools at once.
DO $$
DECLARE
    ref RECORD;
    copy_id UUID;
BEGIN
    FOR ref IN
        SELECT DISTINCT c.grade_id AS old_id, s.id AS school_id
        FROM courses c
        JOIN schools s ON s.created_by = c.teacher_id
        WHERE c.deleted_at IS NULL AND c.grade_id IS NOT NULL
    LOOP
        UPDATE grades SET school_id = ref.school_id
        WHERE id = ref.old_id AND school_id IS NULL;

        IF NOT FOUND THEN
            INSERT INTO grades (name, description, created_by, school_id)
            SELECT g.name, g.description, g.created_by, ref.school_id
            FROM grades g WHERE g.id = ref.old_id
            RETURNING id INTO copy_id;

            UPDATE courses SET grade_id = copy_id
            WHERE grade_id = ref.old_id
              AND deleted_at IS NULL
              AND teacher_id = (SELECT created_by FROM schools WHERE id = ref.school_id);

            -- Grade membership follows the grade: the students in the original
            -- belong to the copy too, or the teacher's students disappear from
            -- their own courses.
            INSERT INTO grade_memberships (grade_id, user_id)
            SELECT copy_id, gm.user_id FROM grade_memberships gm WHERE gm.grade_id = ref.old_id
            ON CONFLICT DO NOTHING;
        END IF;
    END LOOP;

    FOR ref IN
        SELECT DISTINCT c.subject_id AS old_id, s.id AS school_id
        FROM courses c
        JOIN schools s ON s.created_by = c.teacher_id
        WHERE c.deleted_at IS NULL AND c.subject_id IS NOT NULL
    LOOP
        UPDATE subjects SET school_id = ref.school_id
        WHERE id = ref.old_id AND school_id IS NULL;

        IF NOT FOUND THEN
            INSERT INTO subjects (name, description, created_by, school_id)
            SELECT s.name, s.description, s.created_by, ref.school_id
            FROM subjects s WHERE s.id = ref.old_id
            RETURNING id INTO copy_id;

            UPDATE courses SET subject_id = copy_id
            WHERE subject_id = ref.old_id
              AND deleted_at IS NULL
              AND teacher_id = (SELECT created_by FROM schools WHERE id = ref.school_id);
        END IF;
    END LOOP;
END $$;

-- Grades and subjects nobody's courses use go to their creator's school, so
-- nothing is left ownerless.
UPDATE grades g SET school_id = s.id
FROM schools s WHERE s.created_by = g.created_by AND g.school_id IS NULL;

UPDATE subjects sub SET school_id = s.id
FROM schools s WHERE s.created_by = sub.created_by AND sub.school_id IS NULL;

-- Students join the school of every teacher who has them. Both routes count:
-- matching on assignments alone leaves out the students who reach a teacher
-- through one of their courses.
INSERT INTO school_members (school_id, user_id, role)
SELECT DISTINCT s.id, tsa.student_id, 'student'
FROM teacher_student_assignments tsa
JOIN schools s ON s.created_by = tsa.teacher_id
WHERE tsa.status = 'active'
ON CONFLICT DO NOTHING;

INSERT INTO school_members (school_id, user_id, role)
SELECT DISTINCT s.id, e.student_id, 'student'
FROM enrollments e
JOIN courses c ON c.id = e.course_id AND c.deleted_at IS NULL
JOIN schools s ON s.created_by = c.teacher_id
ON CONFLICT DO NOTHING;

INSERT INTO school_members (school_id, user_id, role)
SELECT DISTINCT s.id, gm.user_id, 'student'
FROM grade_memberships gm
JOIN courses c ON c.grade_id = gm.grade_id AND c.deleted_at IS NULL
JOIN schools s ON s.created_by = c.teacher_id
JOIN user_profiles up ON up.id = gm.user_id AND up.profile_type = 'student'
ON CONFLICT DO NOTHING;

-- Names only have to be unique inside a school now.
CREATE UNIQUE INDEX IF NOT EXISTS idx_grades_school_name ON grades(school_id, name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_subjects_school_name ON subjects(school_id, name);

CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    created_by VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE user_profiles
    ADD COLUMN IF NOT EXISTS academic_status VARCHAR(30) NOT NULL DEFAULT 'active'
        CHECK (academic_status IN ('active', 'blocked'));

CREATE TABLE IF NOT EXISTS teacher_student_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(teacher_id, student_id)
);

ALTER TABLE courses
    ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_assignments_teacher ON teacher_student_assignments(teacher_id);
CREATE INDEX IF NOT EXISTS idx_assignments_student ON teacher_student_assignments(student_id);
CREATE INDEX IF NOT EXISTS idx_courses_subject_id ON courses(subject_id);

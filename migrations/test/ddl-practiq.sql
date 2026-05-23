CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_profiles (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(150),
    email VARCHAR(150),
    profile_type VARCHAR(30) NOT NULL DEFAULT 'student' CHECK (profile_type IN ('teacher', 'student')),
    academic_status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (academic_status IN ('active', 'blocked')),
    assistant_base_url TEXT NOT NULL DEFAULT '',
    assistant_api_key TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    created_by VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grade_memberships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grade_id UUID NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(grade_id, user_id)
);

CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    grade_id UUID REFERENCES grades(id) ON DELETE SET NULL,
    subject_id UUID,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    level VARCHAR(50),
    subject VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    created_by VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE courses
    DROP CONSTRAINT IF EXISTS courses_subject_id_fkey;

ALTER TABLE courses
    ADD CONSTRAINT courses_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS teacher_student_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(teacher_id, student_id)
);

CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    status VARCHAR(30) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, student_id)
);

CREATE TABLE IF NOT EXISTS learning_strategies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    code VARCHAR(80) UNIQUE NOT NULL,
    description TEXT,
    status VARCHAR(30) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_learning_strategies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    strategy_id UUID NOT NULL REFERENCES learning_strategies(id) ON DELETE CASCADE,
    is_default BOOLEAN DEFAULT FALSE,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, strategy_id)
);

CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('pdf', 'image', 'video', 'text', 'worksheet')),
    file_url TEXT,
    extracted_text TEXT,
    status VARCHAR(30) DEFAULT 'uploaded',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    order_index INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    material_id UUID REFERENCES materials(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas')),
    question TEXT NOT NULL,
    correct_answer TEXT,
    explanation TEXT,
    difficulty INT DEFAULT 1 CHECK (difficulty BETWEEN 1 AND 10),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS practice_sheets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    strategy_id UUID REFERENCES learning_strategies(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    level INT DEFAULT 1,
    sheet_type VARCHAR(30) NOT NULL DEFAULT 'practice' CHECK (sheet_type IN ('practice', 'level_test')),
    test_style VARCHAR(20) NOT NULL DEFAULT 'keyboard' CHECK (test_style IN ('keyboard', 'canvas')),
    created_by VARCHAR(30) DEFAULT 'teacher' CHECK (created_by IN ('teacher', 'ai')),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS practice_sheet_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    practice_sheet_id UUID NOT NULL REFERENCES practice_sheets(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_index INT DEFAULT 0,
    UNIQUE(practice_sheet_id, exercise_id)
);

CREATE TABLE IF NOT EXISTS student_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    practice_sheet_id UUID REFERENCES practice_sheets(id) ON DELETE SET NULL,
    answer_text TEXT,
    is_correct BOOLEAN,
    score NUMERIC(5,2),
    time_spent_seconds INT DEFAULT 0,
    hints_used INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_work_canvas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES student_attempts(id) ON DELETE CASCADE,
    canvas_data JSONB DEFAULT '{}',
    image_url TEXT,
    recognized_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_topic_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    strategy_id UUID REFERENCES learning_strategies(id) ON DELETE SET NULL,
    mastery_score NUMERIC(5,2) DEFAULT 0,
    current_level INT DEFAULT 1,
    total_attempts INT DEFAULT 0,
    correct_attempts INT DEFAULT 0,
    streak_days INT DEFAULT 0,
    last_practiced_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, topic_id, strategy_id)
);

CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    practice_sheet_id UUID REFERENCES practice_sheets(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    sender VARCHAR(30) NOT NULL CHECK (sender IN ('student', 'ai')),
    message_type VARCHAR(30) DEFAULT 'text' CHECK (message_type IN ('text', 'voice', 'image')),
    content TEXT,
    audio_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_help_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    question TEXT NOT NULL,
    ai_response TEXT,
    help_type VARCHAR(50) CHECK (help_type IN ('hint', 'explanation', 'similar_example')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- NOTEBOOKS
CREATE TABLE IF NOT EXISTS notebooks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id  VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    level       INT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_course_progress (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id    VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id     UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    current_level INT NOT NULL DEFAULT 1,
    updated_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_scp_student_course ON student_course_progress(student_id, course_id);

CREATE TABLE IF NOT EXISTS notebook_pages (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notebook_id  UUID NOT NULL REFERENCES notebooks(id) ON DELETE CASCADE,
    page_number  INT NOT NULL DEFAULT 1,
    title        TEXT DEFAULT '',
    content_type VARCHAR(20) NOT NULL DEFAULT 'canvas',
    content_data TEXT DEFAULT '',
    instructions TEXT DEFAULT '',
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notebook_submissions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id      UUID NOT NULL REFERENCES notebook_pages(id) ON DELETE CASCADE,
    student_id   VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    canvas_data  TEXT DEFAULT '',
    answer_text  TEXT DEFAULT '',
    submitted_at TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    UNIQUE(page_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_notebooks_course ON notebooks(course_id);
CREATE INDEX IF NOT EXISTS idx_notebook_pages_notebook ON notebook_pages(notebook_id);
CREATE INDEX IF NOT EXISTS idx_notebook_submissions_page_student ON notebook_submissions(page_id, student_id);
CREATE INDEX IF NOT EXISTS idx_grade_memberships_user_id ON grade_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_courses_grade_id ON courses(grade_id);
CREATE INDEX IF NOT EXISTS idx_assignments_teacher ON teacher_student_assignments(teacher_id);
CREATE INDEX IF NOT EXISTS idx_assignments_student ON teacher_student_assignments(student_id);
CREATE INDEX IF NOT EXISTS idx_courses_subject_id ON courses(subject_id);

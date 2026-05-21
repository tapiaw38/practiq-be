-- EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USER PROFILES (synced from auth-api-be JWT)
CREATE TABLE IF NOT EXISTS user_profiles (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(150),
    email VARCHAR(150),
    profile_type VARCHAR(30) NOT NULL DEFAULT 'student' CHECK (profile_type IN ('teacher', 'student')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- COURSES
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    level VARCHAR(50),
    subject VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ENROLLMENTS
CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    status VARCHAR(30) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, student_id)
);

-- LEARNING STRATEGIES
CREATE TABLE IF NOT EXISTS learning_strategies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    code VARCHAR(80) UNIQUE NOT NULL,
    description TEXT,
    status VARCHAR(30) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- COURSE LEARNING STRATEGIES
CREATE TABLE IF NOT EXISTS course_learning_strategies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    strategy_id UUID NOT NULL REFERENCES learning_strategies(id) ON DELETE CASCADE,
    is_default BOOLEAN DEFAULT FALSE,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, strategy_id)
);

-- MATERIALS
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

-- TOPICS
CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    order_index INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- EXERCISES
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

-- PRACTICE SHEETS
CREATE TABLE IF NOT EXISTS practice_sheets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    strategy_id UUID REFERENCES learning_strategies(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    level INT DEFAULT 1,
    created_by VARCHAR(30) DEFAULT 'teacher' CHECK (created_by IN ('teacher', 'ai')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- PRACTICE SHEET EXERCISES
CREATE TABLE IF NOT EXISTS practice_sheet_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    practice_sheet_id UUID NOT NULL REFERENCES practice_sheets(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_index INT DEFAULT 0,
    UNIQUE(practice_sheet_id, exercise_id)
);

-- STUDENT ATTEMPTS
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

-- STUDENT WORK CANVAS
CREATE TABLE IF NOT EXISTS student_work_canvas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES student_attempts(id) ON DELETE CASCADE,
    canvas_data JSONB DEFAULT '{}',
    image_url TEXT,
    recognized_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- STUDENT TOPIC PROGRESS
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

-- AI CONVERSATIONS
CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    practice_sheet_id UUID REFERENCES practice_sheets(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI MESSAGES
CREATE TABLE IF NOT EXISTS ai_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    sender VARCHAR(30) NOT NULL CHECK (sender IN ('student', 'ai')),
    message_type VARCHAR(30) DEFAULT 'text' CHECK (message_type IN ('text', 'voice', 'image')),
    content TEXT,
    audio_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI HELP REQUESTS
CREATE TABLE IF NOT EXISTS ai_help_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    question TEXT NOT NULL,
    ai_response TEXT,
    help_type VARCHAR(50) CHECK (help_type IN ('hint', 'explanation', 'similar_example')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Kumon learning strategy
INSERT INTO learning_strategies (name, code, description, status)
VALUES (
    'Kumon Inspired',
    'kumon',
    'Estrategia basada en práctica diaria, repetición, dominio progresivo y avance por niveles.',
    'active'
) ON CONFLICT (code) DO NOTHING;

-- Demo users (these usernames must match auth-api-be usernames)
INSERT INTO user_profiles (id, name, email, profile_type)
VALUES
    ('teacher_demo', 'Profesor Demo', 'teacher@practiq.com', 'teacher'),
    ('student_demo', 'Alumno Demo', 'student@practiq.com', 'student')
ON CONFLICT (id) DO NOTHING;

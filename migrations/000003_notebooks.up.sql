-- Notebooks (cuadernos de tarea creados por el docente, asignados a un curso)
CREATE TABLE IF NOT EXISTS notebooks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id  VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- Pages inside a notebook (each page has teacher content + order)
CREATE TABLE IF NOT EXISTS notebook_pages (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notebook_id  UUID NOT NULL REFERENCES notebooks(id) ON DELETE CASCADE,
    page_number  INT NOT NULL DEFAULT 1,
    title        TEXT DEFAULT '',
    content_type VARCHAR(20) NOT NULL DEFAULT 'canvas', -- 'canvas' | 'text'
    content_data TEXT DEFAULT '',   -- teacher-drawn base64 PNG or HTML text
    instructions TEXT DEFAULT '',   -- optional written instructions above the page
    created_at   TIMESTAMP DEFAULT NOW()
);

-- Student submissions per page
CREATE TABLE IF NOT EXISTS notebook_submissions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id       UUID NOT NULL REFERENCES notebook_pages(id) ON DELETE CASCADE,
    student_id    VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    canvas_data   TEXT DEFAULT '',  -- student's drawing (base64 PNG)
    answer_text   TEXT DEFAULT '',  -- optional text answer
    submitted_at  TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(page_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_notebooks_course ON notebooks(course_id);
CREATE INDEX IF NOT EXISTS idx_notebook_pages_notebook ON notebook_pages(notebook_id);
CREATE INDEX IF NOT EXISTS idx_notebook_submissions_page_student ON notebook_submissions(page_id, student_id);

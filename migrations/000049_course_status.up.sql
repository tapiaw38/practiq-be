-- A course is drafted before anyone sees it, published while it runs, and
-- archived when it is over. Same three states Campus already uses, so the two
-- products do not describe the same lifecycle with different words.
ALTER TABLE courses ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'draft'
    CHECK (status IN ('draft', 'published', 'archived'));

-- Everything that exists today is already being taught, so it is published.
-- The column defaults to 'draft' for what comes next, which is the point:
-- a new course should not reach students the moment it is created.
UPDATE courses SET status = 'published';

CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);

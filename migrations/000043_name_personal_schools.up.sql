-- Auth lives outside this database, but created_by is the teacher username.
-- Rename legacy placeholders immediately instead of waiting for every teacher
-- to sign in and trigger profile sync.
UPDATE schools
SET name = 'Mi escuela de ' || created_by
WHERE kind = 'personal'
  AND name = 'Mi escuela'
  AND COALESCE(created_by, '') <> '';

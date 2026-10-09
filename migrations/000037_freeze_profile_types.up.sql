-- Practiq, not auth-api-be roles, owns academic capability from this point.
-- Existing teacher profiles were already assigned by the previous role bridge;
-- freeze those values and normalize any legacy/empty value to student.
UPDATE user_profiles
SET profile_type = CASE
    WHEN profile_type = 'teacher' THEN 'teacher'
    ELSE 'student'
END;

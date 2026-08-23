-- The flag cleared on the way up cannot be told apart from an answer that was
-- never pending, so this migration restores nothing. Rolling back only brings
-- the previous behaviour, not the old queue contents.
SELECT 1;

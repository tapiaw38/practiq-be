-- Exercises answered by placing options into numbered blanks in the statement.
-- Subject-agnostic: the statement is plain text with {{n}} markers, so it works
-- for prose, formulas or code. The blanks, options and layout live in metadata.
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_type_check;
ALTER TABLE exercises ADD CONSTRAINT exercises_type_check
    CHECK (type IN ('multiple_choice', 'handwritten', 'open_text', 'equation', 'canvas', 'attachment', 'fill_blanks'));

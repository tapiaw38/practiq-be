-- scheduled_at alone described only when a level test opens, and the code
-- treated everything from that instant on as expired, so a scheduled test
-- could never be taken. This is the other end of the window.
--
-- NULL means "open from scheduled_at onwards", which is what the student UI
-- has always promised ("Disponible en la fecha").
ALTER TABLE practice_sheets ADD COLUMN IF NOT EXISTS available_until TIMESTAMP;

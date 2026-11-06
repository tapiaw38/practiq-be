-- A notebook submission is saved after the assistant has looked at it, which
-- takes seconds and happens off the request. Two deliveries from the same
-- student can therefore be in flight at once, and the slow one used to land
-- last and win — overwriting a newer answer with an older one.
--
-- The version is the moment the server accepted the delivery, so ordering
-- comes from arrival rather than from whichever call finished first.
ALTER TABLE notebook_submissions ADD COLUMN IF NOT EXISTS version BIGINT NOT NULL DEFAULT 0;

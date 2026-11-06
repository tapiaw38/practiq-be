-- Only the seed is stored. The avatar is drawn from it on the client, so
-- there is no third-party URL to rot and no free-text field a student could
-- point at an arbitrary image.
ALTER TABLE user_profiles
    ADD COLUMN IF NOT EXISTS avatar_seed VARCHAR(64) NOT NULL DEFAULT '';

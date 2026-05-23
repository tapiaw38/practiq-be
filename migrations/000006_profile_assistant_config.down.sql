ALTER TABLE user_profiles
DROP COLUMN IF EXISTS assistant_base_url,
DROP COLUMN IF EXISTS assistant_api_key;

-- Remove OAuth support
DROP TABLE IF EXISTS oauth_accounts;
DROP INDEX IF EXISTS idx_users_provider_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_password_for_local_users;
ALTER TABLE users DROP COLUMN IF EXISTS provider;
ALTER TABLE users DROP COLUMN IF EXISTS provider_id;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS is_email_verified;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
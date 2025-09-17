-- Remove user profiles and categories tables

DROP TABLE IF EXISTS subscription_history;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS course_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS course_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS user_profiles;

-- Remove optimistic locking column
ALTER TABLE progress DROP COLUMN IF EXISTS version;

-- Remove subscription status enum
DROP TYPE IF EXISTS subscription_status;
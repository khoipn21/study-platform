-- Rollback Enhanced Forum Features Migration

-- Drop new tables
DROP TABLE IF EXISTS forum_topic_subscriptions;
DROP TABLE IF EXISTS forum_notifications;
DROP TABLE IF EXISTS forum_mentions;

-- Remove new columns from forum_posts
ALTER TABLE forum_posts
DROP COLUMN IF EXISTS pin_order,
DROP COLUMN IF EXISTS status;

-- Remove new columns from forum_topics
ALTER TABLE forum_topics
DROP COLUMN IF EXISTS pin_order,
DROP COLUMN IF EXISTS status;
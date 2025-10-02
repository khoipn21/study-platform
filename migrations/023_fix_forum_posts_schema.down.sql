-- Rollback migration 023

BEGIN;

-- Remove added columns
ALTER TABLE forum_posts 
DROP COLUMN IF EXISTS parent_id,
DROP COLUMN IF EXISTS is_edited,
DROP COLUMN IF EXISTS edited_at,
DROP COLUMN IF EXISTS up_votes,
DROP COLUMN IF EXISTS down_votes,
DROP COLUMN IF EXISTS is_answer,
DROP COLUMN IF EXISTS is_pinned;

-- Rename author_id back to user_id
ALTER TABLE forum_posts RENAME COLUMN author_id TO user_id;

-- Restore old indexes
DROP INDEX IF EXISTS idx_forum_posts_author_id;
DROP INDEX IF EXISTS idx_forum_posts_parent_id;
DROP INDEX IF EXISTS idx_forum_posts_is_answer;
DROP INDEX IF EXISTS idx_forum_posts_is_pinned;
DROP INDEX IF EXISTS idx_forum_posts_created_at;
CREATE INDEX IF NOT EXISTS idx_forum_posts_user_id ON forum_posts(user_id);

-- Restore old foreign key
ALTER TABLE forum_posts DROP CONSTRAINT IF EXISTS forum_posts_author_id_fkey;
ALTER TABLE forum_posts ADD CONSTRAINT forum_posts_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add back is_solution column
ALTER TABLE forum_posts ADD COLUMN IF NOT EXISTS is_solution BOOLEAN DEFAULT FALSE;

COMMIT;

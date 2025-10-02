-- Migration to fix forum_posts schema
-- Rename user_id to author_id for consistency
-- Add missing columns

BEGIN;

-- Rename user_id to author_id
ALTER TABLE forum_posts RENAME COLUMN user_id TO author_id;

-- Add missing columns if they don't exist
ALTER TABLE forum_posts 
ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS is_edited BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS edited_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS up_votes INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS down_votes INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS is_answer BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS is_pinned BOOLEAN DEFAULT FALSE;

-- Rename is_solution to is_answer if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'forum_posts' AND column_name = 'is_solution') THEN
        ALTER TABLE forum_posts DROP COLUMN is_solution;
    END IF;
END$$;

-- Update index names
DROP INDEX IF EXISTS idx_forum_posts_user_id;
CREATE INDEX IF NOT EXISTS idx_forum_posts_author_id ON forum_posts(author_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_parent_id ON forum_posts(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_forum_posts_is_answer ON forum_posts(is_answer) WHERE is_answer = TRUE;
CREATE INDEX IF NOT EXISTS idx_forum_posts_is_pinned ON forum_posts(is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX IF NOT EXISTS idx_forum_posts_created_at ON forum_posts(created_at DESC);

-- Update foreign key constraint
ALTER TABLE forum_posts DROP CONSTRAINT IF EXISTS forum_posts_user_id_fkey;
ALTER TABLE forum_posts ADD CONSTRAINT forum_posts_author_id_fkey 
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;

COMMIT;

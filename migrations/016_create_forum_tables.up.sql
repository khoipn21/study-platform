-- Create forum topics table
CREATE TABLE forum_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    created_by_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    is_sticky BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    last_post_at TIMESTAMP NULL,
    last_post_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create forum posts table
CREATE TABLE forum_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMP NULL,
    up_votes INTEGER DEFAULT 0,
    down_votes INTEGER DEFAULT 0,
    is_answer BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create forum votes table
CREATE TABLE forum_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('up', 'down')),
    voted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

-- Create indexes for better performance
CREATE INDEX idx_forum_topics_course_id ON forum_topics(course_id) WHERE course_id IS NOT NULL;
CREATE INDEX idx_forum_topics_created_by ON forum_topics(created_by_id);
CREATE INDEX idx_forum_topics_category ON forum_topics(category);
CREATE INDEX idx_forum_topics_tags ON forum_topics USING GIN(tags);
CREATE INDEX idx_forum_topics_sticky ON forum_topics(is_sticky) WHERE is_sticky = TRUE;
CREATE INDEX idx_forum_topics_updated_at ON forum_topics(updated_at);
CREATE INDEX idx_forum_topics_created_at ON forum_topics(created_at);

CREATE INDEX idx_forum_posts_topic_id ON forum_posts(topic_id);
CREATE INDEX idx_forum_posts_author_id ON forum_posts(author_id);
CREATE INDEX idx_forum_posts_parent_id ON forum_posts(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_forum_posts_created_at ON forum_posts(created_at);
CREATE INDEX idx_forum_posts_votes ON forum_posts(up_votes, down_votes);
CREATE INDEX idx_forum_posts_answer ON forum_posts(is_answer) WHERE is_answer = TRUE;
CREATE INDEX idx_forum_posts_pinned ON forum_posts(is_pinned) WHERE is_pinned = TRUE;

CREATE INDEX idx_forum_votes_post_id ON forum_votes(post_id);
CREATE INDEX idx_forum_votes_user_id ON forum_votes(user_id);
CREATE INDEX idx_forum_votes_type ON forum_votes(vote_type);
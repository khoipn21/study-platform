-- Enhanced Forum Features Migration
-- Add approval system, mentions, pin ordering, and enrollment validation

-- Add status and pin ordering columns to forum_topics
ALTER TABLE forum_topics
ADD COLUMN status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
ADD COLUMN pin_order INTEGER DEFAULT NULL;

-- Add status and pin ordering columns to forum_posts
ALTER TABLE forum_posts
ADD COLUMN status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
ADD COLUMN pin_order INTEGER DEFAULT NULL;

-- Create forum_mentions table for @username functionality
CREATE TABLE forum_mentions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    mentioned_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mentioner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(post_id, mentioned_user_id)
);

-- Create forum_notifications table for mention notifications
CREATE TABLE forum_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('mention', 'topic_approved', 'post_approved', 'topic_reply')),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    reference_id UUID, -- Can reference post_id, topic_id, or mention_id
    reference_type VARCHAR(50) CHECK (reference_type IN ('post', 'topic', 'mention')),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create forum_topic_subscriptions for notification preferences
CREATE TABLE forum_topic_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES forum_topics(id) ON DELETE CASCADE,
    subscribed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, topic_id)
);

-- Add indexes for better performance
CREATE INDEX idx_forum_topics_status ON forum_topics(status);
CREATE INDEX idx_forum_topics_pin_order ON forum_topics(pin_order) WHERE pin_order IS NOT NULL;
CREATE INDEX idx_forum_topics_course_status ON forum_topics(course_id, status) WHERE course_id IS NOT NULL;

CREATE INDEX idx_forum_posts_status ON forum_posts(status);
CREATE INDEX idx_forum_posts_pin_order ON forum_posts(pin_order) WHERE pin_order IS NOT NULL;
CREATE INDEX idx_forum_posts_topic_status ON forum_posts(topic_id, status);

CREATE INDEX idx_forum_mentions_post_id ON forum_mentions(post_id);
CREATE INDEX idx_forum_mentions_mentioned_user ON forum_mentions(mentioned_user_id);
CREATE INDEX idx_forum_mentions_mentioner_user ON forum_mentions(mentioner_user_id);
CREATE INDEX idx_forum_mentions_unread ON forum_mentions(mentioned_user_id, is_read) WHERE is_read = FALSE;

CREATE INDEX idx_forum_notifications_user_id ON forum_notifications(user_id);
CREATE INDEX idx_forum_notifications_unread ON forum_notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_forum_notifications_type ON forum_notifications(type);
CREATE INDEX idx_forum_notifications_reference ON forum_notifications(reference_type, reference_id);

CREATE INDEX idx_forum_subscriptions_user_id ON forum_topic_subscriptions(user_id);
CREATE INDEX idx_forum_subscriptions_topic_id ON forum_topic_subscriptions(topic_id);

-- Add comments for documentation
COMMENT ON COLUMN forum_topics.status IS 'Approval status: pending (needs approval), approved (visible), rejected (hidden)';
COMMENT ON COLUMN forum_topics.pin_order IS 'Order for pinned topics (lower numbers appear first)';
COMMENT ON COLUMN forum_posts.status IS 'Approval status: pending (needs approval), approved (visible), rejected (hidden)';
COMMENT ON COLUMN forum_posts.pin_order IS 'Order for pinned posts within a topic (lower numbers appear first)';
COMMENT ON TABLE forum_mentions IS 'Tracks @username mentions in forum posts';
COMMENT ON TABLE forum_notifications IS 'User notifications for forum events (mentions, approvals, replies)';
COMMENT ON TABLE forum_topic_subscriptions IS 'User subscriptions to forum topics for notifications';
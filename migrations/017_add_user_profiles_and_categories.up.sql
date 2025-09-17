-- Add user profiles and course categories as recommended in architecture review

-- User profiles table for better personalization
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    avatar_url TEXT,
    bio TEXT,
    timezone VARCHAR(50),
    language_preference VARCHAR(10),
    learning_goals TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Categories table for better course organization
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Course categories junction table
CREATE TABLE course_categories (
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, category_id)
);

-- Tags table for flexible course tagging
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Course tags junction table
CREATE TABLE course_tags (
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, tag_id)
);

-- Permissions table for fine-grained access control
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL
);

-- Role permissions junction table
CREATE TABLE role_permissions (
    role VARCHAR(20) NOT NULL,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role, permission_id)
);

-- Progress table optimistic locking
ALTER TABLE progress ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;

-- Subscription management enhancements
CREATE TYPE subscription_status AS ENUM (
    'active', 'past_due', 'canceled', 'unpaid', 'trialing', 'paused'
);

-- Add subscription status to subscriptions table if it exists
-- ALTER TABLE subscriptions ALTER COLUMN status TYPE subscription_status USING status::subscription_status;

-- Subscription history for auditing
CREATE TABLE subscription_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID, -- Will reference subscriptions(id) when table exists
    old_status subscription_status,
    new_status subscription_status,
    reason TEXT,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    changed_by UUID REFERENCES users(id)
);

-- Insert some default permissions
INSERT INTO permissions (name, description, resource, action) VALUES
('courses:create', 'Create new courses', 'courses', 'create'),
('courses:update:own', 'Update own courses', 'courses', 'update'),
('courses:delete:own', 'Delete own courses', 'courses', 'delete'),
('courses:view:all', 'View all courses', 'courses', 'view'),
('forums:moderate', 'Moderate forum posts', 'forums', 'moderate'),
('analytics:view', 'View analytics data', 'analytics', 'view'),
('users:manage', 'Manage user accounts', 'users', 'manage');

-- Insert default role permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'instructor', id FROM permissions WHERE name IN ('courses:create', 'courses:update:own', 'courses:delete:own', 'courses:view:all');

INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions;

-- Add indexes for new tables
CREATE INDEX idx_user_profiles_user ON user_profiles(user_id);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_course_categories_course ON course_categories(course_id);
CREATE INDEX idx_course_categories_category ON course_categories(category_id);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_course_tags_course ON course_tags(course_id);
CREATE INDEX idx_course_tags_tag ON course_tags(tag_id);
CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX idx_subscription_history_subscription ON subscription_history(subscription_id);
CREATE INDEX idx_subscription_history_changed_at ON subscription_history(changed_at DESC);
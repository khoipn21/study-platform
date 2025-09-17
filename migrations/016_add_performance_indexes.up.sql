-- Critical performance indexes as recommended in architecture review

-- Course discovery and filtering indexes
CREATE INDEX IF NOT EXISTS idx_courses_creator ON courses(creator_id);
CREATE INDEX IF NOT EXISTS idx_courses_price ON courses(price, is_free);
CREATE INDEX IF NOT EXISTS idx_courses_created_at ON courses(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_category ON courses(category);
CREATE INDEX IF NOT EXISTS idx_courses_level ON courses(level);

-- Lecture ordering (critical for course display)
CREATE INDEX IF NOT EXISTS idx_lectures_course_sequence ON lectures(course_id, sequence_order);
CREATE INDEX IF NOT EXISTS idx_lectures_status ON lectures(status);

-- Progress tracking (frequent updates/queries)
CREATE INDEX IF NOT EXISTS idx_progress_user ON progress(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_lecture ON progress(lecture_id);
CREATE INDEX IF NOT EXISTS idx_progress_completed ON progress(user_id, completed);
CREATE INDEX IF NOT EXISTS idx_progress_user_course_lecture ON progress(user_id, course_id, lecture_id);

-- Video service performance
CREATE INDEX IF NOT EXISTS idx_videos_lecture ON videos(lecture_id);
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_upload_user ON videos(upload_user_id);

-- Forum performance
CREATE INDEX IF NOT EXISTS idx_forum_topics_course ON forum_topics(course_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic ON forum_posts(topic_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_user ON forum_posts(user_id);
CREATE INDEX IF NOT EXISTS idx_forum_topics_created_at ON forum_topics(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_forum_posts_created_at ON forum_posts(created_at DESC);

-- Payment and enrollment
CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_enrollment_user ON enrollment(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollment_course ON enrollment(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollment_active ON enrollment(is_active, expires_at);

-- Chat history performance
CREATE INDEX IF NOT EXISTS idx_chat_history_user_time ON chat_history(user_id, created_at DESC);

-- Composite indexes for complex queries
CREATE INDEX IF NOT EXISTS idx_courses_search ON courses(is_free, price, created_at DESC) WHERE status = 'published';
CREATE INDEX IF NOT EXISTS idx_enrollment_user_active ON enrollment(user_id, is_active, expires_at) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, created_at DESC);

-- OAuth performance
CREATE INDEX IF NOT EXISTS idx_oauth_accounts_provider ON oauth_accounts(provider, provider_user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_accounts_user ON oauth_accounts(user_id);
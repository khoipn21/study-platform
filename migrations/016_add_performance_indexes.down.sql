-- Remove performance indexes

-- Course indexes
DROP INDEX IF EXISTS idx_courses_creator;
DROP INDEX IF EXISTS idx_courses_price;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_status;
DROP INDEX IF EXISTS idx_courses_category;
DROP INDEX IF EXISTS idx_courses_level;

-- Lecture indexes
DROP INDEX IF EXISTS idx_lectures_course_sequence;
DROP INDEX IF EXISTS idx_lectures_status;

-- Progress indexes
DROP INDEX IF EXISTS idx_progress_user;
DROP INDEX IF EXISTS idx_progress_lecture;
DROP INDEX IF EXISTS idx_progress_completed;
DROP INDEX IF EXISTS idx_progress_user_course_lecture;

-- Video indexes
DROP INDEX IF EXISTS idx_videos_lecture;
DROP INDEX IF EXISTS idx_videos_status;
DROP INDEX IF EXISTS idx_videos_upload_user;

-- Forum indexes
DROP INDEX IF EXISTS idx_forum_topics_course;
DROP INDEX IF EXISTS idx_forum_posts_topic;
DROP INDEX IF EXISTS idx_forum_posts_user;
DROP INDEX IF EXISTS idx_forum_topics_created_at;
DROP INDEX IF EXISTS idx_forum_posts_created_at;

-- Payment and enrollment indexes
DROP INDEX IF EXISTS idx_transactions_user;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_enrollment_user;
DROP INDEX IF EXISTS idx_enrollment_course;
DROP INDEX IF EXISTS idx_enrollment_active;

-- Chat history indexes
DROP INDEX IF EXISTS idx_chat_history_user_time;

-- Composite indexes
DROP INDEX IF EXISTS idx_courses_search;
DROP INDEX IF EXISTS idx_enrollment_user_active;
DROP INDEX IF EXISTS idx_transactions_user_date;

-- OAuth indexes
DROP INDEX IF EXISTS idx_oauth_accounts_provider;
DROP INDEX IF EXISTS idx_oauth_accounts_user;
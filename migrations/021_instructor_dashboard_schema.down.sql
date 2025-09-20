-- Rollback instructor dashboard schema migration

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_instructor_metrics ON transactions;

-- Drop functions
DROP FUNCTION IF EXISTS update_instructor_performance_metrics();
DROP FUNCTION IF EXISTS calculate_video_engagement_score(UUID);

-- Drop views
DROP VIEW IF EXISTS instructor_revenue_analytics;
DROP VIEW IF EXISTS instructor_student_analytics;
DROP VIEW IF EXISTS instructor_course_performance;

-- Drop indexes
DROP INDEX IF EXISTS idx_instructor_dashboard_settings_instructor;
DROP INDEX IF EXISTS idx_course_optimization_suggestions_instructor;
DROP INDEX IF EXISTS idx_course_optimization_suggestions_course;
DROP INDEX IF EXISTS idx_course_optimization_suggestions_status;
DROP INDEX IF EXISTS idx_course_optimization_suggestions_priority;
DROP INDEX IF EXISTS idx_instructor_performance_metrics_instructor_date;
DROP INDEX IF EXISTS idx_instructor_performance_metrics_date;
DROP INDEX IF EXISTS idx_instructor_student_communications_instructor;
DROP INDEX IF EXISTS idx_instructor_student_communications_student;
DROP INDEX IF EXISTS idx_instructor_student_communications_course;
DROP INDEX IF EXISTS idx_instructor_student_communications_type;
DROP INDEX IF EXISTS idx_course_resource_analytics_course_date;
DROP INDEX IF EXISTS idx_course_resource_analytics_resource;
DROP INDEX IF EXISTS idx_course_performance_analytics_course_date;
DROP INDEX IF EXISTS idx_course_performance_analytics_revenue;
DROP INDEX IF EXISTS idx_student_engagement_heatmap_course_lecture;
DROP INDEX IF EXISTS idx_student_engagement_heatmap_timestamp;
DROP INDEX IF EXISTS idx_instructor_team_members_instructor;
DROP INDEX IF EXISTS idx_instructor_team_members_team_member;
DROP INDEX IF EXISTS idx_instructor_team_members_status;
DROP INDEX IF EXISTS idx_instructor_notifications_instructor;
DROP INDEX IF EXISTS idx_instructor_notifications_type;
DROP INDEX IF EXISTS idx_videos_instructor;
DROP INDEX IF EXISTS idx_videos_engagement_score;
DROP INDEX IF EXISTS idx_videos_completion_rate;
DROP INDEX IF EXISTS idx_video_engagement_sessions_video;
DROP INDEX IF EXISTS idx_video_engagement_sessions_user;
DROP INDEX IF EXISTS idx_video_engagement_sessions_engagement;
DROP INDEX IF EXISTS idx_instructor_courses_performance;
DROP INDEX IF EXISTS idx_instructor_revenue_tracking;

-- Remove columns from existing tables
ALTER TABLE courses DROP COLUMN IF EXISTS instructor_notes;
ALTER TABLE courses DROP COLUMN IF EXISTS marketing_description;
ALTER TABLE courses DROP COLUMN IF EXISTS target_audience_description;
ALTER TABLE courses DROP COLUMN IF EXISTS completion_certificate_template;
ALTER TABLE courses DROP COLUMN IF EXISTS auto_approve_enrollments;

ALTER TABLE videos DROP COLUMN IF EXISTS instructor_id;
ALTER TABLE videos DROP COLUMN IF EXISTS engagement_score;
ALTER TABLE videos DROP COLUMN IF EXISTS completion_rate;
ALTER TABLE videos DROP COLUMN IF EXISTS replay_rate;
ALTER TABLE videos DROP COLUMN IF EXISTS ai_questions_count;

-- Drop tables
DROP TABLE IF EXISTS video_engagement_sessions;
DROP TABLE IF EXISTS instructor_notifications;
DROP TABLE IF EXISTS instructor_team_members;
DROP TABLE IF EXISTS student_engagement_heatmap;
DROP TABLE IF EXISTS course_performance_analytics;
DROP TABLE IF EXISTS course_resource_analytics;
DROP TABLE IF EXISTS instructor_student_communications;
DROP TABLE IF EXISTS instructor_performance_metrics;
DROP TABLE IF EXISTS course_optimization_suggestions;
DROP TABLE IF EXISTS instructor_dashboard_settings;
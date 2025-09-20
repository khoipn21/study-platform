-- Instructor Dashboard Migration
-- This migration adds comprehensive instructor dashboard functionality with analytics support

-- Create instructor dashboard settings table
CREATE TABLE IF NOT EXISTS instructor_dashboard_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    dashboard_layout JSONB DEFAULT '{}',
    notification_preferences JSONB DEFAULT '{}',
    analytics_preferences JSONB DEFAULT '{}',
    default_course_settings JSONB DEFAULT '{}',
    ai_assistance_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(instructor_id)
);

-- Create course optimization suggestions from AI
CREATE TABLE IF NOT EXISTS course_optimization_suggestions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    suggestion_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    priority_score DECIMAL(3,2) DEFAULT 0.00,
    expected_impact JSONB DEFAULT '{}',
    implementation_effort VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'pending',
    implemented_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_suggestion_type CHECK (suggestion_type IN (
        'content_gap', 'engagement_improvement', 'pricing_optimization',
        'marketing_enhancement', 'technical_quality', 'accessibility'
    )),
    CONSTRAINT valid_status CHECK (status IN ('pending', 'accepted', 'rejected', 'implemented'))
);

-- Create instructor performance metrics cache
CREATE TABLE IF NOT EXISTS instructor_performance_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    metric_date DATE NOT NULL,
    total_revenue DECIMAL(12,2) DEFAULT 0.00,
    total_enrollments INT DEFAULT 0,
    avg_course_rating DECIMAL(3,2) DEFAULT 0.00,
    total_students INT DEFAULT 0,
    course_completion_rate DECIMAL(5,2) DEFAULT 0.00,
    student_satisfaction_score DECIMAL(3,2) DEFAULT 0.00,
    content_engagement_rate DECIMAL(5,2) DEFAULT 0.00,
    metrics_data JSONB DEFAULT '{}',
    calculated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(instructor_id, metric_date)
);

-- Create student communication history
CREATE TABLE IF NOT EXISTS instructor_student_communications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    communication_type VARCHAR(30) NOT NULL,
    subject VARCHAR(200),
    message TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'sent',
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    read_at TIMESTAMP,
    replied_at TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_communication_type CHECK (communication_type IN (
        'welcome_message', 'milestone_congratulation', 'progress_reminder',
        'course_update', 'direct_message', 'bulk_announcement'
    )),
    CONSTRAINT valid_status CHECK (status IN ('draft', 'scheduled', 'sent', 'delivered', 'failed'))
);

-- Create course resource analytics
CREATE TABLE IF NOT EXISTS course_resource_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    resource_id UUID, -- References files in bucket service
    resource_type VARCHAR(20) NOT NULL,
    download_count INT DEFAULT 0,
    view_count INT DEFAULT 0,
    unique_downloads INT DEFAULT 0,
    unique_views INT DEFAULT 0,
    avg_engagement_time_seconds INT DEFAULT 0,
    last_accessed_at TIMESTAMP,
    analytics_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_resource_type CHECK (resource_type IN (
        'video', 'pdf', 'document', 'image', 'code', 'quiz', 'assignment'
    )),
    UNIQUE(course_id, resource_id, analytics_date)
);

-- Create detailed course performance tracking
CREATE TABLE IF NOT EXISTS course_performance_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    analytics_date DATE NOT NULL,

    -- Discovery metrics
    page_views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    preview_starts INT DEFAULT 0,
    preview_completions INT DEFAULT 0,

    -- Conversion metrics
    enrollment_rate DECIMAL(5,2) DEFAULT 0.00,
    purchase_conversion_rate DECIMAL(5,2) DEFAULT 0.00,
    refund_rate DECIMAL(5,2) DEFAULT 0.00,

    -- Engagement metrics
    avg_video_completion_rate DECIMAL(5,2) DEFAULT 0.00,
    avg_session_duration_seconds INT DEFAULT 0,
    discussion_participation_rate DECIMAL(5,2) DEFAULT 0.00,
    quiz_completion_rate DECIMAL(5,2) DEFAULT 0.00,

    -- Revenue metrics
    gross_revenue DECIMAL(10,2) DEFAULT 0.00,
    net_revenue DECIMAL(10,2) DEFAULT 0.00,
    avg_revenue_per_student DECIMAL(8,2) DEFAULT 0.00,

    -- Quality metrics
    avg_rating DECIMAL(3,2) DEFAULT 0.00,
    review_count INT DEFAULT 0,
    completion_certificate_count INT DEFAULT 0,

    calculated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, analytics_date)
);

-- Create student engagement heatmap data
CREATE TABLE IF NOT EXISTS student_engagement_heatmap (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    lecture_id UUID REFERENCES lectures(id) ON DELETE CASCADE,
    timestamp_seconds INT NOT NULL,
    engagement_score DECIMAL(3,2) DEFAULT 0.00,
    student_count INT DEFAULT 0,
    drop_off_count INT DEFAULT 0,
    replay_count INT DEFAULT 0,
    note_taking_count INT DEFAULT 0,
    quiz_interaction_count INT DEFAULT 0,
    ai_question_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(course_id, lecture_id, timestamp_seconds)
);

-- Create instructor team management table
CREATE TABLE IF NOT EXISTS instructor_team_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    team_member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(30) NOT NULL,
    permissions JSONB DEFAULT '{}',
    course_access UUID[] DEFAULT '{}',
    invited_at TIMESTAMP NOT NULL DEFAULT NOW(),
    joined_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    invitation_token VARCHAR(255),
    invited_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_team_role CHECK (role IN (
        'co_instructor', 'teaching_assistant', 'content_manager', 'marketing_manager'
    )),
    CONSTRAINT valid_team_status CHECK (status IN ('pending', 'active', 'suspended', 'removed')),
    UNIQUE(instructor_id, team_member_id)
);

-- Create instructor notification settings
CREATE TABLE IF NOT EXISTS instructor_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    delivery_method VARCHAR(20) DEFAULT 'email',
    threshold_value DECIMAL(10,2),
    frequency VARCHAR(20) DEFAULT 'immediate',
    last_sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_notification_type CHECK (notification_type IN (
        'new_enrollment', 'course_completion', 'new_review', 'revenue_milestone',
        'engagement_drop', 'technical_issue', 'payment_received', 'refund_issued'
    )),
    CONSTRAINT valid_delivery_method CHECK (delivery_method IN ('email', 'sms', 'push', 'in_app')),
    CONSTRAINT valid_frequency CHECK (frequency IN ('immediate', 'hourly', 'daily', 'weekly')),
    UNIQUE(instructor_id, notification_type, delivery_method)
);

-- Add instructor-specific columns to existing tables
ALTER TABLE courses ADD COLUMN IF NOT EXISTS instructor_notes TEXT;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS marketing_description TEXT;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS target_audience_description TEXT;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS completion_certificate_template UUID;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS auto_approve_enrollments BOOLEAN DEFAULT TRUE;

-- Add video analytics columns for instructor dashboard
ALTER TABLE videos ADD COLUMN IF NOT EXISTS instructor_id UUID REFERENCES users(id);
ALTER TABLE videos ADD COLUMN IF NOT EXISTS engagement_score DECIMAL(3,2) DEFAULT 0.00;
ALTER TABLE videos ADD COLUMN IF NOT EXISTS completion_rate DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE videos ADD COLUMN IF NOT EXISTS replay_rate DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE videos ADD COLUMN IF NOT EXISTS ai_questions_count INT DEFAULT 0;

-- Create enhanced video session tracking for analytics
CREATE TABLE IF NOT EXISTS video_engagement_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_start_time TIMESTAMP NOT NULL DEFAULT NOW(),
    session_end_time TIMESTAMP,
    total_watch_time_seconds INT DEFAULT 0,
    unique_segments_watched INT DEFAULT 0,
    replay_segments_count INT DEFAULT 0,
    pauses_count INT DEFAULT 0,
    seeks_count INT DEFAULT 0,
    quality_changes_count INT DEFAULT 0,
    ai_interactions_count INT DEFAULT 0,
    notes_taken_count INT DEFAULT 0,
    bookmarks_created_count INT DEFAULT 0,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    engagement_score DECIMAL(3,2) DEFAULT 0.00,
    device_type VARCHAR(20),
    browser_type VARCHAR(20),
    connection_quality VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Performance indexes for instructor dashboard
CREATE INDEX IF NOT EXISTS idx_instructor_dashboard_settings_instructor ON instructor_dashboard_settings(instructor_id);

CREATE INDEX IF NOT EXISTS idx_course_optimization_suggestions_instructor ON course_optimization_suggestions(instructor_id);
CREATE INDEX IF NOT EXISTS idx_course_optimization_suggestions_course ON course_optimization_suggestions(course_id);
CREATE INDEX IF NOT EXISTS idx_course_optimization_suggestions_status ON course_optimization_suggestions(status);
CREATE INDEX IF NOT EXISTS idx_course_optimization_suggestions_priority ON course_optimization_suggestions(priority_score DESC);

CREATE INDEX IF NOT EXISTS idx_instructor_performance_metrics_instructor_date ON instructor_performance_metrics(instructor_id, metric_date);
CREATE INDEX IF NOT EXISTS idx_instructor_performance_metrics_date ON instructor_performance_metrics(metric_date);

CREATE INDEX IF NOT EXISTS idx_instructor_student_communications_instructor ON instructor_student_communications(instructor_id);
CREATE INDEX IF NOT EXISTS idx_instructor_student_communications_student ON instructor_student_communications(student_id);
CREATE INDEX IF NOT EXISTS idx_instructor_student_communications_course ON instructor_student_communications(course_id);
CREATE INDEX IF NOT EXISTS idx_instructor_student_communications_type ON instructor_student_communications(communication_type);

CREATE INDEX IF NOT EXISTS idx_course_resource_analytics_course_date ON course_resource_analytics(course_id, analytics_date);
CREATE INDEX IF NOT EXISTS idx_course_resource_analytics_resource ON course_resource_analytics(resource_id);

CREATE INDEX IF NOT EXISTS idx_course_performance_analytics_course_date ON course_performance_analytics(course_id, analytics_date);
CREATE INDEX IF NOT EXISTS idx_course_performance_analytics_revenue ON course_performance_analytics(gross_revenue DESC);

CREATE INDEX IF NOT EXISTS idx_student_engagement_heatmap_course_lecture ON student_engagement_heatmap(course_id, lecture_id);
CREATE INDEX IF NOT EXISTS idx_student_engagement_heatmap_timestamp ON student_engagement_heatmap(timestamp_seconds);

CREATE INDEX IF NOT EXISTS idx_instructor_team_members_instructor ON instructor_team_members(instructor_id);
CREATE INDEX IF NOT EXISTS idx_instructor_team_members_team_member ON instructor_team_members(team_member_id);
CREATE INDEX IF NOT EXISTS idx_instructor_team_members_status ON instructor_team_members(status);

CREATE INDEX IF NOT EXISTS idx_instructor_notifications_instructor ON instructor_notifications(instructor_id);
CREATE INDEX IF NOT EXISTS idx_instructor_notifications_type ON instructor_notifications(notification_type);

CREATE INDEX IF NOT EXISTS idx_videos_instructor ON videos(instructor_id);
CREATE INDEX IF NOT EXISTS idx_videos_engagement_score ON videos(engagement_score DESC);
CREATE INDEX IF NOT EXISTS idx_videos_completion_rate ON videos(completion_rate DESC);

CREATE INDEX IF NOT EXISTS idx_video_engagement_sessions_video ON video_engagement_sessions(video_id);
CREATE INDEX IF NOT EXISTS idx_video_engagement_sessions_user ON video_engagement_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_video_engagement_sessions_engagement ON video_engagement_sessions(engagement_score DESC);

-- Create composite indexes for common dashboard queries
CREATE INDEX IF NOT EXISTS idx_instructor_courses_performance ON courses(creator_id, status, average_rating DESC, total_enrollments DESC);
CREATE INDEX IF NOT EXISTS idx_instructor_revenue_tracking ON transactions(user_id, status, created_at DESC) WHERE status = 'completed';

-- Create views for instructor dashboard analytics
CREATE OR REPLACE VIEW instructor_revenue_analytics AS
SELECT
    c.creator_id as instructor_id,
    DATE_TRUNC('day', t.created_at) as revenue_date,
    COUNT(t.id) as total_transactions,
    SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END) as daily_revenue,
    COUNT(DISTINCT t.user_id) as unique_customers,
    COUNT(DISTINCT t.course_id) as courses_sold,
    AVG(CASE WHEN t.status = 'completed' THEN t.amount END) as avg_transaction_value
FROM courses c
JOIN transactions t ON c.id = t.course_id
WHERE t.status IN ('completed', 'refunded')
GROUP BY c.creator_id, DATE_TRUNC('day', t.created_at)
ORDER BY revenue_date DESC;

CREATE OR REPLACE VIEW instructor_student_analytics AS
SELECT
    c.creator_id as instructor_id,
    COUNT(DISTINCT e.user_id) as total_students,
    COUNT(DISTINCT CASE WHEN e.payment_status = 'paid' THEN e.user_id END) as paying_students,
    COUNT(DISTINCT p.user_id) as active_students,
    AVG(p.completion_percentage) as avg_completion_rate,
    COUNT(DISTINCT c.id) as total_courses
FROM courses c
LEFT JOIN enrollments e ON c.id = e.course_id
LEFT JOIN progress p ON c.id = p.course_id AND p.completion_percentage > 0
GROUP BY c.creator_id;

CREATE OR REPLACE VIEW instructor_course_performance AS
SELECT
    c.id as course_id,
    c.creator_id as instructor_id,
    c.title,
    c.average_rating,
    c.total_enrollments,
    COUNT(DISTINCT e.user_id) as total_students,
    COUNT(DISTINCT CASE WHEN e.payment_status = 'paid' THEN e.user_id END) as paying_students,
    AVG(p.completion_percentage) as avg_completion_rate,
    SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END) as total_revenue,
    COUNT(DISTINCT t.id) as total_sales
FROM courses c
LEFT JOIN enrollments e ON c.id = e.course_id
LEFT JOIN progress p ON c.id = p.course_id
LEFT JOIN transactions t ON c.id = t.course_id AND t.status = 'completed'
GROUP BY c.id, c.creator_id, c.title, c.average_rating, c.total_enrollments
ORDER BY total_revenue DESC;

-- Create triggers for automatic metric calculation
CREATE OR REPLACE FUNCTION update_instructor_performance_metrics() RETURNS TRIGGER AS $$
BEGIN
    -- Update daily metrics when transactions are completed
    IF NEW.status = 'completed' AND OLD.status != 'completed' THEN
        INSERT INTO instructor_performance_metrics (
            instructor_id,
            metric_date,
            total_revenue,
            total_enrollments
        )
        SELECT
            c.creator_id,
            CURRENT_DATE,
            COALESCE(SUM(t.amount), 0),
            COUNT(DISTINCT e.user_id)
        FROM courses c
        LEFT JOIN transactions t ON c.id = t.course_id AND t.status = 'completed' AND DATE(t.created_at) = CURRENT_DATE
        LEFT JOIN enrollments e ON c.id = e.course_id AND DATE(e.created_at) = CURRENT_DATE
        WHERE c.creator_id = (SELECT creator_id FROM courses WHERE id = NEW.course_id)
        GROUP BY c.creator_id
        ON CONFLICT (instructor_id, metric_date)
        DO UPDATE SET
            total_revenue = EXCLUDED.total_revenue,
            total_enrollments = EXCLUDED.total_enrollments,
            calculated_at = NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for instructor metrics updates
DROP TRIGGER IF EXISTS trigger_update_instructor_metrics ON transactions;
CREATE TRIGGER trigger_update_instructor_metrics
    AFTER UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_instructor_performance_metrics();

-- Function to calculate video engagement scores
CREATE OR REPLACE FUNCTION calculate_video_engagement_score(video_uuid UUID) RETURNS DECIMAL(3,2) AS $$
DECLARE
    engagement_score DECIMAL(3,2);
BEGIN
    SELECT
        CASE
            WHEN COUNT(*) = 0 THEN 0.0
            ELSE LEAST(1.0, GREATEST(0.0,
                (AVG(completion_percentage) * 0.4) +
                (AVG(CASE WHEN replay_segments_count > 0 THEN 0.3 ELSE 0.0 END) * 0.2) +
                (AVG(CASE WHEN ai_interactions_count > 0 THEN 0.2 ELSE 0.0 END) * 0.2) +
                (AVG(CASE WHEN notes_taken_count > 0 THEN 0.1 ELSE 0.0 END) * 0.2)
            ))
        END INTO engagement_score
    FROM video_engagement_sessions
    WHERE video_id = video_uuid;

    RETURN COALESCE(engagement_score, 0.0);
END;
$$ LANGUAGE plpgsql;

-- Add default instructor dashboard settings for existing instructors
INSERT INTO instructor_dashboard_settings (instructor_id, dashboard_layout, notification_preferences, analytics_preferences)
SELECT
    id,
    '{"widgets": ["revenue", "students", "courses", "analytics"], "layout": "grid"}',
    '{"email": true, "push": true, "new_enrollments": true, "revenue_milestones": true}',
    '{"show_advanced": false, "default_timeframe": "30_days", "include_free_courses": true}'
FROM users
WHERE role = 'instructor'
ON CONFLICT (instructor_id) DO NOTHING;

-- Add instructor_id to existing videos based on course creator
UPDATE videos SET instructor_id = (
    SELECT creator_id FROM courses WHERE courses.id = videos.course_id
)
WHERE instructor_id IS NULL AND course_id IS NOT NULL;

-- Comments for documentation
COMMENT ON TABLE instructor_dashboard_settings IS 'Instructor dashboard configuration and preferences';
COMMENT ON TABLE course_optimization_suggestions IS 'AI-generated suggestions for course improvement';
COMMENT ON TABLE instructor_performance_metrics IS 'Daily aggregated performance metrics for instructors';
COMMENT ON TABLE instructor_student_communications IS 'Communication history between instructors and students';
COMMENT ON TABLE course_resource_analytics IS 'Analytics for course resources like videos, documents, etc.';
COMMENT ON TABLE course_performance_analytics IS 'Detailed course performance metrics';
COMMENT ON TABLE student_engagement_heatmap IS 'Granular engagement data for video content';
COMMENT ON TABLE instructor_team_members IS 'Team management for instructors with multiple collaborators';
COMMENT ON TABLE instructor_notifications IS 'Notification preferences for instructors';
COMMENT ON TABLE video_engagement_sessions IS 'Detailed video viewing sessions for analytics';

COMMENT ON VIEW instructor_revenue_analytics IS 'Daily revenue analytics aggregated by instructor';
COMMENT ON VIEW instructor_student_analytics IS 'Student engagement and enrollment analytics';
COMMENT ON VIEW instructor_course_performance IS 'Comprehensive course performance metrics';
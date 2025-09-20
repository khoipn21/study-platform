-- Enhanced Video Streaming Platform Schema
-- This migration creates comprehensive tables for production-grade video streaming
-- with network quality monitoring, analytics, and real-time adaptation capabilities

-- Update existing videos table to support advanced video streaming features
ALTER TABLE videos DROP COLUMN IF EXISTS bucket_path;
ALTER TABLE videos DROP COLUMN IF EXISTS formats;

-- Add new columns to videos table for Cloudflare Stream integration
ALTER TABLE videos
ADD COLUMN IF NOT EXISTS cloudflare_uid VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS upload_user_id UUID NOT NULL DEFAULT gen_random_uuid(),
ADD COLUMN IF NOT EXISTS course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'processing',
ADD COLUMN IF NOT EXISTS visibility VARCHAR(20) DEFAULT 'private',
ADD COLUMN IF NOT EXISTS stream_url TEXT,
ADD COLUMN IF NOT EXISTS preview_url TEXT,
ADD COLUMN IF NOT EXISTS file_size_bytes BIGINT,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Create video_qualities table for adaptive streaming support
CREATE TABLE IF NOT EXISTS video_qualities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    quality_label VARCHAR(20) NOT NULL, -- '360p', '720p', '1080p', '1440p', '2160p'
    bitrate_kbps INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    fps INTEGER DEFAULT 30,
    codec VARCHAR(20) DEFAULT 'h264',
    url TEXT NOT NULL,
    file_size_bytes BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create viewing_sessions table for tracking user video sessions
CREATE TABLE IF NOT EXISTS viewing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_heartbeat TIMESTAMP NOT NULL DEFAULT NOW(),
    current_time_seconds INTEGER DEFAULT 0,
    current_quality VARCHAR(20),
    total_watch_time_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    ip_address INET,
    device_type VARCHAR(50),
    screen_resolution VARCHAR(20),
    connection_type VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create network_metrics table for tracking network quality data
CREATE TABLE IF NOT EXISTS network_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    bandwidth_mbps DECIMAL(10,2),
    latency_ms INTEGER,
    packet_loss_percent DECIMAL(5,2),
    connection_type VARCHAR(20), -- 'wifi', '4g', '5g', 'ethernet', 'unknown'
    quality_score INTEGER, -- 1-10 scale
    recommended_quality VARCHAR(20),
    buffer_health_seconds INTEGER,
    jitter_ms INTEGER,
    throughput_mbps DECIMAL(10,2),
    signal_strength INTEGER, -- for mobile connections
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create video_analytics table for comprehensive video performance tracking
CREATE TABLE IF NOT EXISTS video_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_views INTEGER DEFAULT 0,
    unique_viewers INTEGER DEFAULT 0,
    total_watch_time_seconds BIGINT DEFAULT 0,
    avg_watch_time_seconds INTEGER DEFAULT 0,
    completion_rate DECIMAL(5,2) DEFAULT 0.00,
    quality_distribution JSONB DEFAULT '{}', -- {"360p": 30, "720p": 50, "1080p": 20}
    geographic_distribution JSONB DEFAULT '{}', -- {"US": 60, "EU": 30, "AS": 10}
    device_distribution JSONB DEFAULT '{}', -- {"desktop": 70, "mobile": 25, "tablet": 5}
    peak_concurrent_viewers INTEGER DEFAULT 0,
    total_buffering_events INTEGER DEFAULT 0,
    avg_startup_time_seconds DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(video_id, date)
);

-- Create quality_changes table for tracking adaptive streaming decisions
CREATE TABLE IF NOT EXISTS quality_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    from_quality VARCHAR(20),
    to_quality VARCHAR(20) NOT NULL,
    reason VARCHAR(100), -- 'bandwidth_drop', 'bandwidth_increase', 'buffer_low', 'user_manual'
    trigger_metric VARCHAR(50), -- 'bandwidth', 'latency', 'packet_loss', 'buffer_health'
    metric_value DECIMAL(10,2),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    auto_switch BOOLEAN DEFAULT TRUE
);

-- Create video_permissions table for access control
CREATE TABLE IF NOT EXISTS video_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID,
    role_id UUID,
    permission_type VARCHAR(20) NOT NULL, -- 'view', 'download', 'share', 'edit'
    granted_by UUID NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (user_id IS NOT NULL OR role_id IS NOT NULL)
);

-- Create video_engagement_events table for detailed user interaction tracking
CREATE TABLE IF NOT EXISTS video_engagement_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    event_type VARCHAR(30) NOT NULL, -- 'play', 'pause', 'seek', 'quality_change', 'fullscreen', 'volume_change'
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    video_timestamp_seconds INTEGER,
    event_data JSONB DEFAULT '{}', -- Additional event-specific data
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create buffering_events table for tracking video buffering performance
CREATE TABLE IF NOT EXISTS buffering_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_ms INTEGER,
    video_timestamp_seconds INTEGER,
    quality VARCHAR(20),
    buffer_before_seconds DECIMAL(5,2),
    buffer_after_seconds DECIMAL(5,2),
    network_bandwidth_mbps DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create error_logs table for tracking video streaming errors
CREATE TABLE IF NOT EXISTS video_error_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID,
    error_type VARCHAR(50) NOT NULL, -- 'network_error', 'decode_error', 'format_error', 'permission_error'
    error_code VARCHAR(20),
    error_message TEXT,
    video_timestamp_seconds INTEGER,
    quality VARCHAR(20),
    user_agent TEXT,
    network_info JSONB DEFAULT '{}',
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Performance Indexes for optimal query performance

-- Videos table indexes
CREATE INDEX IF NOT EXISTS idx_videos_cloudflare_uid ON videos(cloudflare_uid);
CREATE INDEX IF NOT EXISTS idx_videos_upload_user_id ON videos(upload_user_id);
CREATE INDEX IF NOT EXISTS idx_videos_course_id ON videos(course_id);
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_visibility ON videos(visibility);
CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at);
CREATE INDEX IF NOT EXISTS idx_videos_deleted_at ON videos(deleted_at) WHERE deleted_at IS NULL;

-- Video qualities indexes
CREATE INDEX IF NOT EXISTS idx_video_qualities_video_id ON video_qualities(video_id);
CREATE INDEX IF NOT EXISTS idx_video_qualities_bitrate ON video_qualities(bitrate_kbps);

-- Viewing sessions indexes
CREATE INDEX IF NOT EXISTS idx_viewing_sessions_session_id ON viewing_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_viewing_sessions_user_id ON viewing_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_viewing_sessions_video_id ON viewing_sessions(video_id);
CREATE INDEX IF NOT EXISTS idx_viewing_sessions_started_at ON viewing_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_viewing_sessions_last_heartbeat ON viewing_sessions(last_heartbeat);

-- Network metrics indexes for time-series queries
CREATE INDEX IF NOT EXISTS idx_network_metrics_session_id ON network_metrics(session_id);
CREATE INDEX IF NOT EXISTS idx_network_metrics_timestamp ON network_metrics(timestamp);
CREATE INDEX IF NOT EXISTS idx_network_metrics_user_video ON network_metrics(user_id, video_id);
CREATE INDEX IF NOT EXISTS idx_network_metrics_quality_score ON network_metrics(quality_score);

-- Video analytics indexes
CREATE INDEX IF NOT EXISTS idx_video_analytics_video_date ON video_analytics(video_id, date);
CREATE INDEX IF NOT EXISTS idx_video_analytics_date ON video_analytics(date);

-- Quality changes indexes
CREATE INDEX IF NOT EXISTS idx_quality_changes_session_id ON quality_changes(session_id);
CREATE INDEX IF NOT EXISTS idx_quality_changes_timestamp ON quality_changes(timestamp);
CREATE INDEX IF NOT EXISTS idx_quality_changes_reason ON quality_changes(reason);

-- Video permissions indexes
CREATE INDEX IF NOT EXISTS idx_video_permissions_video_id ON video_permissions(video_id);
CREATE INDEX IF NOT EXISTS idx_video_permissions_user_id ON video_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_video_permissions_expires_at ON video_permissions(expires_at) WHERE expires_at IS NOT NULL;

-- Engagement events indexes
CREATE INDEX IF NOT EXISTS idx_engagement_events_session_id ON video_engagement_events(session_id);
CREATE INDEX IF NOT EXISTS idx_engagement_events_event_type ON video_engagement_events(event_type);
CREATE INDEX IF NOT EXISTS idx_engagement_events_timestamp ON video_engagement_events(timestamp);

-- Buffering events indexes
CREATE INDEX IF NOT EXISTS idx_buffering_events_session_id ON buffering_events(session_id);
CREATE INDEX IF NOT EXISTS idx_buffering_events_started_at ON buffering_events(started_at);

-- Error logs indexes
CREATE INDEX IF NOT EXISTS idx_video_error_logs_session_id ON video_error_logs(session_id);
CREATE INDEX IF NOT EXISTS idx_video_error_logs_error_type ON video_error_logs(error_type);
CREATE INDEX IF NOT EXISTS idx_video_error_logs_timestamp ON video_error_logs(timestamp);

-- Foreign key constraints
ALTER TABLE videos
ADD CONSTRAINT fk_videos_upload_user_id
FOREIGN KEY (upload_user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE viewing_sessions
ADD CONSTRAINT fk_viewing_sessions_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE network_metrics
ADD CONSTRAINT fk_network_metrics_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE quality_changes
ADD CONSTRAINT fk_quality_changes_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE video_permissions
ADD CONSTRAINT fk_video_permissions_granted_by
FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE video_engagement_events
ADD CONSTRAINT fk_video_engagement_events_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE buffering_events
ADD CONSTRAINT fk_buffering_events_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create enum types for better data integrity
DO $$ BEGIN
    CREATE TYPE video_status AS ENUM ('uploading', 'processing', 'ready', 'error', 'deleted');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE video_visibility AS ENUM ('public', 'unlisted', 'private', 'course_only');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE connection_type AS ENUM ('wifi', '4g', '5g', 'ethernet', 'unknown');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add check constraints for data validation
ALTER TABLE video_qualities
ADD CONSTRAINT chk_video_qualities_bitrate CHECK (bitrate_kbps > 0),
ADD CONSTRAINT chk_video_qualities_dimensions CHECK (width > 0 AND height > 0),
ADD CONSTRAINT chk_video_qualities_fps CHECK (fps > 0 AND fps <= 120);

ALTER TABLE network_metrics
ADD CONSTRAINT chk_network_metrics_bandwidth CHECK (bandwidth_mbps >= 0),
ADD CONSTRAINT chk_network_metrics_latency CHECK (latency_ms >= 0),
ADD CONSTRAINT chk_network_metrics_packet_loss CHECK (packet_loss_percent >= 0 AND packet_loss_percent <= 100),
ADD CONSTRAINT chk_network_metrics_quality_score CHECK (quality_score >= 1 AND quality_score <= 10);

ALTER TABLE video_analytics
ADD CONSTRAINT chk_video_analytics_completion_rate CHECK (completion_rate >= 0 AND completion_rate <= 100);

-- Add triggers for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_videos_updated_at BEFORE UPDATE ON videos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_video_analytics_updated_at BEFORE UPDATE ON video_analytics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create materialized views for performance analytics
CREATE MATERIALIZED VIEW IF NOT EXISTS video_performance_summary AS
SELECT
    v.id as video_id,
    v.title,
    v.course_id,
    COUNT(DISTINCT vs.user_id) as total_unique_viewers,
    COUNT(vs.id) as total_sessions,
    AVG(vs.total_watch_time_seconds) as avg_watch_time,
    AVG(nm.quality_score) as avg_quality_score,
    COUNT(DISTINCT qc.id) as total_quality_changes,
    COUNT(DISTINCT be.id) as total_buffering_events,
    MAX(vs.last_heartbeat) as last_viewed_at
FROM videos v
LEFT JOIN viewing_sessions vs ON v.id = vs.video_id
LEFT JOIN network_metrics nm ON v.id = nm.video_id
LEFT JOIN quality_changes qc ON v.id = qc.video_id
LEFT JOIN buffering_events be ON v.id = be.video_id
WHERE v.deleted_at IS NULL
GROUP BY v.id, v.title, v.course_id;

-- Create unique indexes for materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_video_performance_summary_video_id
ON video_performance_summary(video_id);

COMMENT ON TABLE videos IS 'Enhanced video metadata table with Cloudflare Stream integration and comprehensive tracking';
COMMENT ON TABLE video_qualities IS 'Available quality variants for each video supporting adaptive streaming';
COMMENT ON TABLE viewing_sessions IS 'Active and historical video viewing sessions with real-time tracking';
COMMENT ON TABLE network_metrics IS 'Time-series network quality metrics for adaptive streaming decisions';
COMMENT ON TABLE video_analytics IS 'Aggregated daily analytics for video performance and engagement';
COMMENT ON TABLE quality_changes IS 'Log of all quality changes during video playback for analysis';
COMMENT ON TABLE video_permissions IS 'Access control permissions for videos at user and role level';
COMMENT ON TABLE video_engagement_events IS 'Detailed user interaction events during video playback';
COMMENT ON TABLE buffering_events IS 'Buffering incidents tracking for performance optimization';
COMMENT ON TABLE video_error_logs IS 'Error tracking for debugging and system monitoring';
-- Drop enhanced video streaming platform schema
-- This down migration removes all enhancements for video streaming platform

-- Drop materialized views
DROP MATERIALIZED VIEW IF EXISTS video_performance_summary;

-- Drop triggers
DROP TRIGGER IF EXISTS update_videos_updated_at ON videos;
DROP TRIGGER IF EXISTS update_video_analytics_updated_at ON video_analytics;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS video_error_logs CASCADE;
DROP TABLE IF EXISTS buffering_events CASCADE;
DROP TABLE IF EXISTS video_engagement_events CASCADE;
DROP TABLE IF EXISTS video_permissions CASCADE;
DROP TABLE IF EXISTS quality_changes CASCADE;
DROP TABLE IF EXISTS video_analytics CASCADE;
DROP TABLE IF EXISTS network_metrics CASCADE;
DROP TABLE IF EXISTS viewing_sessions CASCADE;
DROP TABLE IF EXISTS video_qualities CASCADE;

-- Remove columns added to videos table
ALTER TABLE videos
DROP COLUMN IF EXISTS cloudflare_uid,
DROP COLUMN IF EXISTS upload_user_id,
DROP COLUMN IF EXISTS course_id,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS visibility,
DROP COLUMN IF EXISTS stream_url,
DROP COLUMN IF EXISTS preview_url,
DROP COLUMN IF EXISTS file_size_bytes,
DROP COLUMN IF EXISTS metadata,
DROP COLUMN IF EXISTS deleted_at;

-- Restore original videos table columns
ALTER TABLE videos
ADD COLUMN IF NOT EXISTS bucket_path TEXT NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS formats JSONB DEFAULT '{}';

-- Drop enum types
DROP TYPE IF EXISTS video_status CASCADE;
DROP TYPE IF EXISTS video_visibility CASCADE;
DROP TYPE IF EXISTS connection_type CASCADE;
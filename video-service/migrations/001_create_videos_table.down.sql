-- Drop indexes
DROP INDEX IF EXISTS idx_videos_deleted_at;
DROP INDEX IF EXISTS idx_videos_status;
DROP INDEX IF EXISTS idx_videos_course_lecture;
DROP INDEX IF EXISTS idx_videos_user_id;
DROP INDEX IF EXISTS idx_videos_cloudflare_uid;

-- Drop videos table
DROP TABLE IF EXISTS videos;
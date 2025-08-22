-- Drop indexes
DROP INDEX IF EXISTS idx_viewing_sessions_session;
DROP INDEX IF EXISTS idx_viewing_sessions_user_video;

-- Drop viewing_sessions table
DROP TABLE IF EXISTS viewing_sessions;
-- User viewing sessions
CREATE TABLE viewing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    started_at TIMESTAMP DEFAULT NOW(),
    last_heartbeat TIMESTAMP DEFAULT NOW(),
    current_time_seconds INTEGER DEFAULT 0,
    current_quality VARCHAR(10),
    total_watch_time_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    ip_address INET,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for viewing_sessions table
CREATE INDEX idx_viewing_sessions_user_video ON viewing_sessions(user_id, video_id);
CREATE INDEX idx_viewing_sessions_session ON viewing_sessions(session_id);
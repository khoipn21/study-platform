-- Video access permissions
CREATE TABLE video_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID,
    role_id UUID,
    permission_type VARCHAR(20) NOT NULL, -- 'view', 'download', 'share'
    granted_by UUID NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Analytics aggregations
CREATE TABLE video_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_views INTEGER DEFAULT 0,
    unique_viewers INTEGER DEFAULT 0,
    total_watch_time_seconds BIGINT DEFAULT 0,
    avg_watch_time_seconds INTEGER DEFAULT 0,
    completion_rate DECIMAL(5,2) DEFAULT 0,
    quality_distribution JSONB, -- {'360p': 30, '720p': 50, '1080p': 20}
    geographic_distribution JSONB,
    device_distribution JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(video_id, date)
);

-- Indexes for video_permissions and video_analytics tables
CREATE INDEX idx_video_permissions_video_user ON video_permissions(video_id, user_id);
CREATE INDEX idx_video_analytics_video_date ON video_analytics(video_id, date);
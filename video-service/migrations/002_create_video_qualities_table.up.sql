-- Video quality variants
CREATE TABLE video_qualities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    quality_label VARCHAR(10) NOT NULL, -- '360p', '720p', '1080p'
    bitrate_kbps INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    fps INTEGER DEFAULT 30,
    codec VARCHAR(20) DEFAULT 'h264',
    url VARCHAR(500) NOT NULL,
    file_size_bytes BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for video_qualities table
CREATE INDEX idx_video_qualities_video_id ON video_qualities(video_id);
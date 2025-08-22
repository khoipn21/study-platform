-- Video metadata table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cloudflare_uid VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_seconds INTEGER,
    file_size_bytes BIGINT,
    upload_user_id UUID NOT NULL,
    course_id UUID,
    lecture_id UUID,
    status VARCHAR(20) DEFAULT 'processing', -- 'uploading', 'processing', 'ready', 'error'
    visibility VARCHAR(20) DEFAULT 'private', -- 'public', 'private', 'unlisted'
    thumbnail_url VARCHAR(500),
    stream_url VARCHAR(500),
    preview_url VARCHAR(500),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Indexes for videos table
CREATE INDEX idx_videos_cloudflare_uid ON videos(cloudflare_uid);
CREATE INDEX idx_videos_user_id ON videos(upload_user_id);
CREATE INDEX idx_videos_course_lecture ON videos(course_id, lecture_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_deleted_at ON videos(deleted_at);
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecture_id UUID REFERENCES lectures(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    bucket_path TEXT NOT NULL,
    formats JSONB,
    duration INT,
    thumbnail_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_videos_lecture_id ON videos(lecture_id);
CREATE INDEX idx_videos_bucket_path ON videos(bucket_path);
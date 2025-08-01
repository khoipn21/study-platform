-- Create upload sessions table for multipart uploads
CREATE TABLE upload_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    upload_id VARCHAR(255) NOT NULL, -- S3 multipart upload ID
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    total_size BIGINT NOT NULL,
    uploaded_size BIGINT DEFAULT 0,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'aborted')),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_upload_sessions_user_status ON upload_sessions(user_id, status);
CREATE INDEX idx_upload_sessions_expires_at ON upload_sessions(expires_at);
CREATE INDEX idx_upload_sessions_upload_id ON upload_sessions(upload_id);
CREATE INDEX idx_upload_sessions_bucket_key ON upload_sessions(bucket_name, object_key);

-- Create unique constraint on upload_id
CREATE UNIQUE INDEX idx_upload_sessions_upload_id_unique ON upload_sessions(upload_id);
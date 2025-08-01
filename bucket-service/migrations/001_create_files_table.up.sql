-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create files table
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    upload_user_id UUID NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    checksum VARCHAR(64),
    thumbnail_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX idx_files_user_id ON files(upload_user_id);
CREATE INDEX idx_files_bucket_key ON files(bucket_name, object_key);
CREATE INDEX idx_files_created_at ON files(created_at);
CREATE INDEX idx_files_deleted_at ON files(deleted_at);
CREATE INDEX idx_files_content_type ON files(content_type);
CREATE INDEX idx_files_is_public ON files(is_public);

-- Create unique constraint on bucket_name + object_key for active files
CREATE UNIQUE INDEX idx_files_bucket_object_active ON files(bucket_name, object_key) WHERE deleted_at IS NULL;
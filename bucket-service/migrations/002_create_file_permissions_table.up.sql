-- Create file permissions table
CREATE TABLE file_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    user_id UUID,
    permission_type VARCHAR(20) NOT NULL CHECK (permission_type IN ('read', 'write', 'delete')),
    granted_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_file_permissions_file_user ON file_permissions(file_id, user_id);
CREATE INDEX idx_file_permissions_user ON file_permissions(user_id);
CREATE INDEX idx_file_permissions_type ON file_permissions(permission_type);

-- Prevent duplicate permissions for the same file/user/type combination
CREATE UNIQUE INDEX idx_file_permissions_unique ON file_permissions(file_id, user_id, permission_type);
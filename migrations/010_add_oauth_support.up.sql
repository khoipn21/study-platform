-- Add OAuth support to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50) DEFAULT 'local';
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider_id VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_email_verified BOOLEAN DEFAULT FALSE;

-- Make password_hash nullable for OAuth users
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Create OAuth accounts table for linking multiple providers
CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(100) NOT NULL,
    provider_email VARCHAR(100),
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

-- Create indexes for OAuth accounts
CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);
CREATE INDEX idx_oauth_accounts_provider ON oauth_accounts(provider);
CREATE INDEX idx_oauth_accounts_provider_id ON oauth_accounts(provider_id);

-- Update users table constraints
CREATE UNIQUE INDEX idx_users_provider_id ON users(provider, provider_id) WHERE provider != 'local';

-- Add constraint: local users must have password, OAuth users don't need it
ALTER TABLE users ADD CONSTRAINT check_password_for_local_users 
    CHECK (
        (provider = 'local' AND password_hash IS NOT NULL) OR 
        (provider != 'local' AND password_hash IS NULL) OR
        (provider != 'local' AND password_hash IS NOT NULL)
    );
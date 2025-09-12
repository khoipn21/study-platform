-- ==============================================
-- STUDY PLATFORM - SECURE DATABASE INITIALIZATION
-- ==============================================
-- This script sets up the database with proper security configurations
-- and creates necessary users with minimal required privileges
-- ==============================================

-- Create application users with specific privileges
DO $$
BEGIN
    -- Check if the application user already exists
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = current_setting('app.db_user', true)) THEN
        -- Create application user with limited privileges
        EXECUTE format('CREATE USER %I WITH PASSWORD %L', 
                      current_setting('app.db_user', true), 
                      current_setting('app.db_password', true));
    END IF;
END
$$;

-- Create read-only user for analytics/reporting
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = current_setting('app.readonly_user', true)) THEN
        EXECUTE format('CREATE USER %I WITH PASSWORD %L', 
                      current_setting('app.readonly_user', true), 
                      current_setting('app.readonly_password', true));
    END IF;
END
$$;

-- Grant database-level privileges
GRANT CONNECT ON DATABASE studyplatform TO current_setting('app.db_user', true);
GRANT CONNECT ON DATABASE studyplatform TO current_setting('app.readonly_user', true);

-- Create application schema
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION current_setting('app.db_user', true);
CREATE SCHEMA IF NOT EXISTS analytics AUTHORIZATION current_setting('app.readonly_user', true);

-- Set secure default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO current_setting('app.db_user', true);
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT USAGE, SELECT ON SEQUENCES TO current_setting('app.db_user', true);
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON TABLES TO current_setting('app.readonly_user', true);

-- Security configuration
SET session_replication_role = replica;

-- Enable row level security by default
ALTER SYSTEM SET row_security = on;

-- Set secure password settings
ALTER SYSTEM SET password_encryption = 'scram-sha-256';

-- Configure connection security
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_ciphers = 'ECDHE+AESGCM:ECDHE+CHACHA20:DHE+AESGCM:DHE+CHACHA20:!aNULL:!MD5:!DSS';
ALTER SYSTEM SET ssl_prefer_server_ciphers = on;
ALTER SYSTEM SET ssl_min_protocol_version = 'TLSv1.2';

-- Set logging for security monitoring
ALTER SYSTEM SET log_connections = on;
ALTER SYSTEM SET log_disconnections = on;
ALTER SYSTEM SET log_checkpoints = on;
ALTER SYSTEM SET log_lock_waits = on;
ALTER SYSTEM SET log_statement = 'ddl';
ALTER SYSTEM SET log_min_duration_statement = 1000;

-- Reload configuration
SELECT pg_reload_conf();

-- Create security audit table
CREATE TABLE IF NOT EXISTS app.security_audit (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    user_id INTEGER,
    ip_address INET,
    user_agent TEXT,
    request_path TEXT,
    http_method VARCHAR(10),
    status_code INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    details JSONB
);

-- Create index on audit table for performance
CREATE INDEX IF NOT EXISTS idx_security_audit_created_at ON app.security_audit(created_at);
CREATE INDEX IF NOT EXISTS idx_security_audit_event_type ON app.security_audit(event_type);
CREATE INDEX IF NOT EXISTS idx_security_audit_user_id ON app.security_audit(user_id);
CREATE INDEX IF NOT EXISTS idx_security_audit_ip_address ON app.security_audit(ip_address);

-- Create function to log security events
CREATE OR REPLACE FUNCTION app.log_security_event(
    p_event_type VARCHAR(50),
    p_user_id INTEGER DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL,
    p_request_path TEXT DEFAULT NULL,
    p_http_method VARCHAR(10) DEFAULT NULL,
    p_status_code INTEGER DEFAULT NULL,
    p_details JSONB DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    INSERT INTO app.security_audit (
        event_type, user_id, ip_address, user_agent, 
        request_path, http_method, status_code, details
    ) VALUES (
        p_event_type, p_user_id, p_ip_address, p_user_agent,
        p_request_path, p_http_method, p_status_code, p_details
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to clean old audit logs
CREATE OR REPLACE FUNCTION app.cleanup_security_audit(retention_days INTEGER DEFAULT 90) 
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM app.security_audit 
    WHERE created_at < NOW() - INTERVAL '1 day' * retention_days;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create users table with security features
CREATE TABLE IF NOT EXISTS app.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'student' CHECK (role IN ('student', 'instructor', 'admin')),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    password_changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    mfa_enabled BOOLEAN DEFAULT false,
    mfa_secret VARCHAR(32),
    backup_codes TEXT[],
    session_token_hash VARCHAR(255),
    session_expires_at TIMESTAMP WITH TIME ZONE
);

-- Create secure indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON app.users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON app.users(username);
CREATE INDEX IF NOT EXISTS idx_users_session_token ON app.users(session_token_hash) WHERE session_token_hash IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_locked_until ON app.users(locked_until) WHERE locked_until IS NOT NULL;

-- Enable Row Level Security on users table
ALTER TABLE app.users ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY users_own_data ON app.users
    FOR ALL TO current_setting('app.db_user', true)
    USING (id = current_setting('app.current_user_id', true)::INTEGER);

-- Create function to hash passwords securely
CREATE OR REPLACE FUNCTION app.hash_password(password TEXT, salt TEXT DEFAULT NULL)
RETURNS TABLE(hash TEXT, salt_out TEXT) AS $$
DECLARE
    generated_salt TEXT;
BEGIN
    -- Generate salt if not provided
    IF salt IS NULL THEN
        generated_salt := encode(gen_random_bytes(16), 'hex');
    ELSE
        generated_salt := salt;
    END IF;
    
    -- Return hash and salt
    RETURN QUERY SELECT 
        crypt(password, '$2b$12$' || generated_salt) as hash,
        generated_salt as salt_out;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to verify passwords
CREATE OR REPLACE FUNCTION app.verify_password(password TEXT, hash TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN hash = crypt(password, hash);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to handle failed login attempts
CREATE OR REPLACE FUNCTION app.handle_failed_login(user_email TEXT)
RETURNS VOID AS $$
DECLARE
    max_attempts INTEGER := 5;
    lockout_duration INTERVAL := '15 minutes';
BEGIN
    UPDATE app.users 
    SET 
        failed_login_attempts = failed_login_attempts + 1,
        locked_until = CASE 
            WHEN failed_login_attempts + 1 >= max_attempts 
            THEN NOW() + lockout_duration 
            ELSE locked_until 
        END
    WHERE email = user_email;
    
    -- Log security event
    PERFORM app.log_security_event(
        'FAILED_LOGIN_ATTEMPT',
        (SELECT id FROM app.users WHERE email = user_email),
        NULL, -- IP will be set by application
        NULL, -- User agent will be set by application
        '/api/v1/auth/login',
        'POST',
        401,
        jsonb_build_object('email', user_email)
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to handle successful login
CREATE OR REPLACE FUNCTION app.handle_successful_login(user_email TEXT, login_ip INET)
RETURNS VOID AS $$
BEGIN
    UPDATE app.users 
    SET 
        failed_login_attempts = 0,
        locked_until = NULL,
        last_login_at = NOW(),
        last_login_ip = login_ip
    WHERE email = user_email;
    
    -- Log security event
    PERFORM app.log_security_event(
        'SUCCESSFUL_LOGIN',
        (SELECT id FROM app.users WHERE email = user_email),
        login_ip,
        NULL, -- User agent will be set by application
        '/api/v1/auth/login',
        'POST',
        200,
        jsonb_build_object('email', user_email)
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION app.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to users table
DROP TRIGGER IF EXISTS trigger_users_updated_at ON app.users;
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON app.users
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

-- Create view for safe user data (no sensitive fields)
CREATE OR REPLACE VIEW app.users_safe AS
SELECT 
    id,
    username,
    email,
    role,
    is_active,
    email_verified,
    created_at,
    updated_at,
    last_login_at,
    mfa_enabled
FROM app.users;

-- Grant appropriate permissions
GRANT SELECT ON app.users_safe TO current_setting('app.db_user', true);
GRANT SELECT ON app.users_safe TO current_setting('app.readonly_user', true);

-- Create session management table
CREATE TABLE IF NOT EXISTS app.user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT true
);

-- Create indexes for session management
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON app.user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON app.user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON app.user_sessions(expires_at);

-- Function to cleanup expired sessions
CREATE OR REPLACE FUNCTION app.cleanup_expired_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM app.user_sessions WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    -- Log cleanup event
    PERFORM app.log_security_event(
        'SESSION_CLEANUP',
        NULL,
        NULL,
        NULL,
        '/internal/cleanup',
        'SYSTEM',
        200,
        jsonb_build_object('deleted_sessions', deleted_count)
    );
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Grant execute permissions on functions
GRANT EXECUTE ON FUNCTION app.hash_password(TEXT, TEXT) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.verify_password(TEXT, TEXT) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.handle_failed_login(TEXT) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.handle_successful_login(TEXT, INET) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.log_security_event(VARCHAR, INTEGER, INET, TEXT, TEXT, VARCHAR, INTEGER, JSONB) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.cleanup_security_audit(INTEGER) TO current_setting('app.db_user', true);
GRANT EXECUTE ON FUNCTION app.cleanup_expired_sessions() TO current_setting('app.db_user', true);

-- Final security settings
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO current_setting('app.db_user', true);
GRANT USAGE ON SCHEMA public TO current_setting('app.readonly_user', true);

-- Log successful initialization
INSERT INTO app.security_audit (event_type, details) 
VALUES ('DATABASE_INIT', jsonb_build_object('timestamp', NOW(), 'version', '1.0.0'));

COMMIT;
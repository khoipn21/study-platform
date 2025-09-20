-- Enhanced Payment Integration Migration
-- This migration adds comprehensive payment tracking, access control, and audit capabilities

-- Add payment verification columns to enrollments
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) DEFAULT 'pending'
    CHECK (payment_status IN ('pending', 'paid', 'refunded', 'expired'));
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS payment_verified_at TIMESTAMP;
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS access_expires_at TIMESTAMP;
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS transaction_id UUID REFERENCES transactions(id);

-- Create course access logs table for audit trail
CREATE TABLE IF NOT EXISTS course_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    course_id UUID NOT NULL REFERENCES courses(id),
    lecture_id UUID REFERENCES lectures(id),
    access_type VARCHAR(20) NOT NULL CHECK (access_type IN ('full', 'preview', 'denied')),
    access_granted BOOLEAN NOT NULL,
    payment_required BOOLEAN NOT NULL,
    payment_verified BOOLEAN DEFAULT FALSE,
    transaction_id UUID REFERENCES transactions(id),
    client_ip INET,
    user_agent TEXT,
    access_duration_seconds INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create access cache table for performance optimization
CREATE TABLE IF NOT EXISTS course_access_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    access_level VARCHAR(20) NOT NULL CHECK (access_level IN ('full', 'preview', 'denied')),
    payment_verified BOOLEAN NOT NULL,
    transaction_id UUID,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- Create lecture preview sessions table
CREATE TABLE IF NOT EXISTS lecture_preview_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    lecture_id UUID NOT NULL REFERENCES lectures(id),
    session_started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    session_duration_seconds INTEGER NOT NULL DEFAULT 0,
    preview_limit_seconds INTEGER NOT NULL DEFAULT 600,
    preview_exhausted BOOLEAN DEFAULT FALSE,
    ip_address INET,
    last_accessed_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, lecture_id)
);

-- Create audit logs table for comprehensive tracking
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    course_id UUID REFERENCES courses(id),
    lecture_id UUID REFERENCES lectures(id),
    transaction_id UUID REFERENCES transactions(id),
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Enhanced transactions table with additional payment tracking
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(50) DEFAULT 'lemonsqueezy';
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_verified_at TIMESTAMP;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS refunded_at TIMESTAMP;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_method_details JSONB;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS risk_score DECIMAL(3,2) DEFAULT 0.0;

-- Enhanced courses table for better payment integration
ALTER TABLE courses ADD COLUMN IF NOT EXISTS preview_enabled BOOLEAN DEFAULT TRUE;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS preview_duration_minutes INTEGER DEFAULT 10;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS requires_enrollment_approval BOOLEAN DEFAULT FALSE;

-- Enhanced lectures table for granular access control
ALTER TABLE lectures ADD COLUMN IF NOT EXISTS preview_available BOOLEAN DEFAULT FALSE;
ALTER TABLE lectures ADD COLUMN IF NOT EXISTS access_level VARCHAR(20) DEFAULT 'paid'
    CHECK (access_level IN ('free', 'paid', 'preview'));

-- Create payment events table for webhook tracking
CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(50) NOT NULL,
    provider VARCHAR(20) NOT NULL DEFAULT 'lemonsqueezy',
    provider_event_id VARCHAR(100) NOT NULL,
    transaction_id UUID REFERENCES transactions(id),
    user_id UUID REFERENCES users(id),
    course_id UUID REFERENCES courses(id),
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_event_id)
);

-- Create user payment methods table for stored payment info
CREATE TABLE IF NOT EXISTS user_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    provider VARCHAR(20) NOT NULL DEFAULT 'lemonsqueezy',
    provider_customer_id VARCHAR(100),
    payment_method_type VARCHAR(20) NOT NULL DEFAULT 'card',
    last_four_digits VARCHAR(4),
    expiry_month INTEGER,
    expiry_year INTEGER,
    brand VARCHAR(20),
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_course_access_logs_user_course ON course_access_logs(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_course_access_logs_user_lecture ON course_access_logs(user_id, lecture_id);
CREATE INDEX IF NOT EXISTS idx_course_access_logs_created_at ON course_access_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_course_access_logs_access_type ON course_access_logs(access_type);

CREATE INDEX IF NOT EXISTS idx_course_access_cache_user_course ON course_access_cache(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_course_access_cache_expires_at ON course_access_cache(expires_at);
CREATE INDEX IF NOT EXISTS idx_course_access_cache_access_level ON course_access_cache(access_level);

CREATE INDEX IF NOT EXISTS idx_lecture_preview_sessions_user_lecture ON lecture_preview_sessions(user_id, lecture_id);
CREATE INDEX IF NOT EXISTS idx_lecture_preview_sessions_exhausted ON lecture_preview_sessions(preview_exhausted);
CREATE INDEX IF NOT EXISTS idx_lecture_preview_sessions_last_accessed ON lecture_preview_sessions(last_accessed_at);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_course_id ON audit_logs(course_id);

CREATE INDEX IF NOT EXISTS idx_enrollments_payment_status ON enrollments(payment_status);
CREATE INDEX IF NOT EXISTS idx_enrollments_payment_verified_at ON enrollments(payment_verified_at);
CREATE INDEX IF NOT EXISTS idx_enrollments_access_expires_at ON enrollments(access_expires_at);
CREATE INDEX IF NOT EXISTS idx_enrollments_transaction_id ON enrollments(transaction_id);

CREATE INDEX IF NOT EXISTS idx_transactions_payment_provider ON transactions(payment_provider);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_verified_at ON transactions(payment_verified_at);
CREATE INDEX IF NOT EXISTS idx_transactions_user_course ON transactions(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status_updated ON transactions(status, updated_at);

CREATE INDEX IF NOT EXISTS idx_payment_events_type_processed ON payment_events(event_type, processed);
CREATE INDEX IF NOT EXISTS idx_payment_events_provider_event_id ON payment_events(provider, provider_event_id);
CREATE INDEX IF NOT EXISTS idx_payment_events_transaction_id ON payment_events(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payment_events_user_course ON payment_events(user_id, course_id);

CREATE INDEX IF NOT EXISTS idx_user_payment_methods_user_id ON user_payment_methods(user_id);
CREATE INDEX IF NOT EXISTS idx_user_payment_methods_default ON user_payment_methods(user_id, is_default) WHERE is_default = TRUE;
CREATE INDEX IF NOT EXISTS idx_user_payment_methods_active ON user_payment_methods(user_id, is_active) WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_courses_preview_enabled ON courses(preview_enabled);
CREATE INDEX IF NOT EXISTS idx_lectures_access_level ON lectures(access_level);
CREATE INDEX IF NOT EXISTS idx_lectures_preview_available ON lectures(preview_available);

-- Create composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_course_access_user_course_time ON course_access_logs(user_id, course_id, created_at);
CREATE INDEX IF NOT EXISTS idx_enrollment_user_payment_status ON enrollments(user_id, payment_status);
CREATE INDEX IF NOT EXISTS idx_transaction_user_status_course ON transactions(user_id, status, course_id);

-- Create partial indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transactions_pending ON transactions(user_id, course_id, created_at)
    WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_transactions_completed ON transactions(user_id, course_id, payment_verified_at)
    WHERE status = 'completed';

-- Create function for automatic cache expiration cleanup
CREATE OR REPLACE FUNCTION cleanup_expired_access_cache() RETURNS void AS $$
BEGIN
    DELETE FROM course_access_cache WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Create function for updating enrollment payment status
CREATE OR REPLACE FUNCTION update_enrollment_payment_status() RETURNS TRIGGER AS $$
BEGIN
    -- When a transaction is completed, update corresponding enrollment
    IF NEW.status = 'completed' AND OLD.status != 'completed' AND NEW.course_id IS NOT NULL THEN
        UPDATE enrollments
        SET payment_status = 'paid',
            payment_verified_at = NEW.payment_verified_at,
            transaction_id = NEW.id,
            updated_at = NOW()
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;

        -- Clear access cache for this user/course combination
        DELETE FROM course_access_cache
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;
    END IF;

    -- When a transaction is refunded, update enrollment status
    IF NEW.status = 'refunded' AND OLD.status != 'refunded' AND NEW.course_id IS NOT NULL THEN
        UPDATE enrollments
        SET payment_status = 'refunded',
            updated_at = NOW()
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;

        -- Clear access cache
        DELETE FROM course_access_cache
        WHERE user_id = NEW.user_id AND course_id = NEW.course_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for automatic enrollment updates
DROP TRIGGER IF EXISTS trigger_update_enrollment_payment_status ON transactions;
CREATE TRIGGER trigger_update_enrollment_payment_status
    AFTER UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_enrollment_payment_status();

-- Create function for automatic preview session cleanup
CREATE OR REPLACE FUNCTION cleanup_old_preview_sessions() RETURNS void AS $$
BEGIN
    -- Mark sessions as exhausted if they've exceeded their time limit
    UPDATE lecture_preview_sessions
    SET preview_exhausted = TRUE,
        updated_at = NOW()
    WHERE preview_exhausted = FALSE
      AND (EXTRACT(EPOCH FROM (NOW() - session_started_at)) > preview_limit_seconds);

    -- Clean up old exhausted sessions (older than 30 days)
    DELETE FROM lecture_preview_sessions
    WHERE preview_exhausted = TRUE
      AND created_at < NOW() - INTERVAL '30 days';
END;
$$ LANGUAGE plpgsql;

-- Add constraints for data integrity
ALTER TABLE course_access_cache ADD CONSTRAINT check_cache_expires_future
    CHECK (expires_at > created_at);

ALTER TABLE lecture_preview_sessions ADD CONSTRAINT check_preview_limit_positive
    CHECK (preview_limit_seconds > 0);

ALTER TABLE lecture_preview_sessions ADD CONSTRAINT check_session_duration_valid
    CHECK (session_duration_seconds >= 0 AND session_duration_seconds <= preview_limit_seconds);

-- Update existing data to set defaults
UPDATE enrollments SET payment_status = 'paid'
WHERE status = 'enrolled' AND payment_status IS NULL;

UPDATE courses SET preview_enabled = TRUE, preview_duration_minutes = 10
WHERE preview_enabled IS NULL;

UPDATE lectures SET access_level = CASE
    WHEN is_free = TRUE THEN 'free'
    ELSE 'paid'
END
WHERE access_level IS NULL;

-- Set first lecture of each course as preview available
UPDATE lectures SET preview_available = TRUE
WHERE id IN (
    SELECT DISTINCT ON (course_id) id
    FROM lectures
    WHERE access_level = 'paid'
    ORDER BY course_id, order_number
);

-- Create view for payment analytics
CREATE OR REPLACE VIEW payment_analytics AS
SELECT
    DATE_TRUNC('day', t.created_at) as payment_date,
    COUNT(*) as total_transactions,
    COUNT(CASE WHEN t.status = 'completed' THEN 1 END) as successful_payments,
    COUNT(CASE WHEN t.status = 'failed' THEN 1 END) as failed_payments,
    COUNT(CASE WHEN t.status = 'refunded' THEN 1 END) as refunded_payments,
    SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END) as total_revenue,
    AVG(CASE WHEN t.status = 'completed' THEN t.amount END) as avg_transaction_amount,
    COUNT(DISTINCT t.user_id) as unique_customers,
    COUNT(DISTINCT t.course_id) as courses_purchased
FROM transactions t
WHERE t.course_id IS NOT NULL
GROUP BY DATE_TRUNC('day', t.created_at)
ORDER BY payment_date DESC;

-- Create view for course access analytics
CREATE OR REPLACE VIEW course_access_analytics AS
SELECT
    c.id as course_id,
    c.title as course_title,
    c.is_paid,
    COUNT(DISTINCT cal.user_id) as total_access_attempts,
    COUNT(DISTINCT CASE WHEN cal.access_granted = TRUE THEN cal.user_id END) as successful_accesses,
    COUNT(DISTINCT CASE WHEN cal.access_type = 'preview' THEN cal.user_id END) as preview_accesses,
    COUNT(DISTINCT e.user_id) as total_enrollments,
    COUNT(DISTINCT CASE WHEN e.payment_status = 'paid' THEN e.user_id END) as paid_enrollments,
    AVG(cal.access_duration_seconds) as avg_access_duration,
    SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END) as total_revenue
FROM courses c
LEFT JOIN course_access_logs cal ON c.id = cal.course_id
LEFT JOIN enrollments e ON c.id = e.course_id
LEFT JOIN transactions t ON c.id = t.course_id AND t.status = 'completed'
GROUP BY c.id, c.title, c.is_paid
ORDER BY total_revenue DESC;

-- Create indexes on the views for better performance
CREATE INDEX IF NOT EXISTS idx_payment_analytics_date ON transactions(DATE_TRUNC('day', created_at));
CREATE INDEX IF NOT EXISTS idx_course_access_analytics_course ON course_access_logs(course_id, access_granted, access_type);

-- Grant necessary permissions (adjust as needed for your setup)
-- GRANT SELECT ON payment_analytics TO readonly_role;
-- GRANT SELECT ON course_access_analytics TO readonly_role;

COMMENT ON TABLE course_access_logs IS 'Comprehensive audit trail of all course access attempts';
COMMENT ON TABLE course_access_cache IS 'Performance cache for course access validation results';
COMMENT ON TABLE lecture_preview_sessions IS 'Tracking table for lecture preview sessions and time limits';
COMMENT ON TABLE audit_logs IS 'General audit logging for payment and access events';
COMMENT ON TABLE payment_events IS 'Webhook events and payment provider communications';
COMMENT ON TABLE user_payment_methods IS 'Stored payment methods for users';

COMMENT ON COLUMN enrollments.payment_status IS 'Current payment status: pending, paid, refunded, expired';
COMMENT ON COLUMN enrollments.payment_verified_at IS 'Timestamp when payment was verified by provider';
COMMENT ON COLUMN enrollments.access_expires_at IS 'When access expires (for subscription-based courses)';
COMMENT ON COLUMN transactions.payment_verified_at IS 'When payment was verified by webhook or API';
COMMENT ON COLUMN transactions.risk_score IS 'Fraud risk score from 0.0 (safe) to 1.0 (high risk)';
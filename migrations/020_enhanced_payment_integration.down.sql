-- Rollback Enhanced Payment Integration Migration

-- Drop views
DROP VIEW IF EXISTS payment_analytics;
DROP VIEW IF EXISTS course_access_analytics;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_enrollment_payment_status ON transactions;

-- Drop functions
DROP FUNCTION IF EXISTS update_enrollment_payment_status();
DROP FUNCTION IF EXISTS cleanup_expired_access_cache();
DROP FUNCTION IF EXISTS cleanup_old_preview_sessions();

-- Drop indexes
DROP INDEX IF EXISTS idx_payment_analytics_date;
DROP INDEX IF EXISTS idx_course_access_analytics_course;

DROP INDEX IF EXISTS idx_course_access_logs_user_course;
DROP INDEX IF EXISTS idx_course_access_logs_user_lecture;
DROP INDEX IF EXISTS idx_course_access_logs_created_at;
DROP INDEX IF EXISTS idx_course_access_logs_access_type;

DROP INDEX IF EXISTS idx_course_access_cache_user_course;
DROP INDEX IF EXISTS idx_course_access_cache_expires_at;
DROP INDEX IF EXISTS idx_course_access_cache_access_level;

DROP INDEX IF EXISTS idx_lecture_preview_sessions_user_lecture;
DROP INDEX IF EXISTS idx_lecture_preview_sessions_exhausted;
DROP INDEX IF EXISTS idx_lecture_preview_sessions_last_accessed;

DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_course_id;

DROP INDEX IF EXISTS idx_enrollments_payment_status;
DROP INDEX IF EXISTS idx_enrollments_payment_verified_at;
DROP INDEX IF EXISTS idx_enrollments_access_expires_at;
DROP INDEX IF EXISTS idx_enrollments_transaction_id;

DROP INDEX IF EXISTS idx_transactions_payment_provider;
DROP INDEX IF EXISTS idx_transactions_payment_verified_at;
DROP INDEX IF EXISTS idx_transactions_user_course;
DROP INDEX IF EXISTS idx_transactions_status_updated;

DROP INDEX IF EXISTS idx_payment_events_type_processed;
DROP INDEX IF EXISTS idx_payment_events_provider_event_id;
DROP INDEX IF EXISTS idx_payment_events_transaction_id;
DROP INDEX IF EXISTS idx_payment_events_user_course;

DROP INDEX IF EXISTS idx_user_payment_methods_user_id;
DROP INDEX IF EXISTS idx_user_payment_methods_default;
DROP INDEX IF EXISTS idx_user_payment_methods_active;

DROP INDEX IF EXISTS idx_courses_preview_enabled;
DROP INDEX IF EXISTS idx_lectures_access_level;
DROP INDEX IF EXISTS idx_lectures_preview_available;

DROP INDEX IF EXISTS idx_course_access_user_course_time;
DROP INDEX IF EXISTS idx_enrollment_user_payment_status;
DROP INDEX IF EXISTS idx_transaction_user_status_course;

DROP INDEX IF EXISTS idx_transactions_pending;
DROP INDEX IF EXISTS idx_transactions_completed;

-- Drop new tables
DROP TABLE IF EXISTS course_access_logs;
DROP TABLE IF EXISTS course_access_cache;
DROP TABLE IF EXISTS lecture_preview_sessions;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS user_payment_methods;

-- Remove columns from existing tables
ALTER TABLE enrollments DROP COLUMN IF EXISTS payment_status;
ALTER TABLE enrollments DROP COLUMN IF EXISTS payment_verified_at;
ALTER TABLE enrollments DROP COLUMN IF EXISTS access_expires_at;
ALTER TABLE enrollments DROP COLUMN IF EXISTS transaction_id;

ALTER TABLE transactions DROP COLUMN IF EXISTS payment_provider;
ALTER TABLE transactions DROP COLUMN IF EXISTS payment_verified_at;
ALTER TABLE transactions DROP COLUMN IF EXISTS refunded_at;
ALTER TABLE transactions DROP COLUMN IF EXISTS expires_at;
ALTER TABLE transactions DROP COLUMN IF EXISTS payment_method_details;
ALTER TABLE transactions DROP COLUMN IF EXISTS risk_score;

ALTER TABLE courses DROP COLUMN IF EXISTS preview_enabled;
ALTER TABLE courses DROP COLUMN IF EXISTS preview_duration_minutes;
ALTER TABLE courses DROP COLUMN IF EXISTS requires_enrollment_approval;

ALTER TABLE lectures DROP COLUMN IF EXISTS preview_available;
ALTER TABLE lectures DROP COLUMN IF EXISTS access_level;
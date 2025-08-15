-- Master seed script for Study Platform development data (corrected for actual schema)
-- Run this script to populate the database with development data

-- Note: This script assumes all migrations have been run successfully

BEGIN;

-- Disable triggers and foreign key checks for faster inserts
SET session_replication_role = 'replica';

-- Clear existing data (in reverse dependency order)
-- Only truncate tables that actually exist
TRUNCATE TABLE forum_posts CASCADE;
TRUNCATE TABLE forum_topics CASCADE;
TRUNCATE TABLE chat_history CASCADE;
TRUNCATE TABLE subscriptions CASCADE;
TRUNCATE TABLE transactions CASCADE;
TRUNCATE TABLE payment_methods CASCADE;
TRUNCATE TABLE enrollments CASCADE; 
TRUNCATE TABLE progress CASCADE;
TRUNCATE TABLE lectures CASCADE;
TRUNCATE TABLE courses CASCADE;
TRUNCATE TABLE oauth_accounts CASCADE;
TRUNCATE TABLE users CASCADE;

-- Seed data files (in dependency order)
\i 001_seed_users.sql
\i 002_seed_courses.sql  
\i 003_seed_enrollments.sql
\i 004_seed_payment_data.sql
\i 005_seed_forum_data_fixed.sql
\i 006_seed_chat_data_corrected.sql

-- Re-enable triggers and foreign key checks
SET session_replication_role = 'origin';

-- Update statistics
ANALYZE;

COMMIT;

-- Display summary of seeded data
SELECT 
    'Users' as table_name, 
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT role::text, ', ') as details
FROM users
UNION ALL
SELECT 
    'Courses' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT level::text, ', ') as details
FROM courses  
UNION ALL
SELECT 
    'Lectures' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status::text, ', ') as details
FROM lectures
UNION ALL
SELECT 
    'Enrollments' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status::text, ', ') as details
FROM enrollments
UNION ALL
SELECT 
    'Payment Methods' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT provider, ', ') as details
FROM payment_methods
UNION ALL
SELECT 
    'Transactions' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status, ', ') as details
FROM transactions
UNION ALL
SELECT 
    'Subscriptions' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status, ', ') as details
FROM subscriptions
UNION ALL
SELECT 
    'Forum Topics' as table_name,
    COUNT(*) as record_count,
    COUNT(CASE WHEN is_pinned THEN 1 END)::text || ' pinned' as details
FROM forum_topics
UNION ALL
SELECT 
    'Forum Posts' as table_name,
    COUNT(*) as record_count,
    COUNT(CASE WHEN is_solution THEN 1 END)::text || ' solutions' as details
FROM forum_posts
UNION ALL
SELECT 
    'Chat Messages' as table_name,
    COUNT(*) as record_count,
    COUNT(CASE WHEN is_user THEN 1 END)::text || ' user, ' || 
    COUNT(CASE WHEN NOT is_user THEN 1 END)::text || ' bot' as details
FROM chat_history;

-- Display some sample data for verification
\echo
\echo 'Sample Users:'
SELECT username, email, role FROM users LIMIT 5;

\echo
\echo 'Sample Courses:'
SELECT title, instructor_name, level, price FROM courses WHERE status = 'published' LIMIT 3;

\echo  
\echo 'Sample Enrollments:'
SELECT u.username, c.title, e.status, e.progress_percentage 
FROM enrollments e 
JOIN users u ON e.user_id = u.id 
JOIN courses c ON e.course_id = c.id 
LIMIT 5;
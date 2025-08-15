-- Master seed script for Study Platform development data
-- Run this script to populate the database with development data

-- Note: This script assumes all migrations have been run successfully

BEGIN;

-- Disable triggers and foreign key checks for faster inserts
SET session_replication_role = 'replica';

-- Clear existing data (in reverse dependency order)
TRUNCATE TABLE forum_votes CASCADE;
TRUNCATE TABLE forum_posts CASCADE;
TRUNCATE TABLE forum_topics CASCADE;
TRUNCATE TABLE chat_history CASCADE;
TRUNCATE TABLE chat_sessions CASCADE;
TRUNCATE TABLE subscriptions CASCADE;
TRUNCATE TABLE transactions CASCADE;
TRUNCATE TABLE payment_methods CASCADE;
TRUNCATE TABLE enrollments CASCADE;
TRUNCATE TABLE lectures CASCADE;
TRUNCATE TABLE courses CASCADE;
TRUNCATE TABLE oauth_accounts CASCADE;
TRUNCATE TABLE users CASCADE;

-- Reset sequences
SELECT setval('users_id_seq', 1, false) WHERE EXISTS (SELECT 1 FROM pg_class WHERE relname = 'users_id_seq');

-- Seed data files (in dependency order)
\i 001_seed_users.sql
\i 002_seed_courses.sql
\i 003_seed_enrollments.sql
\i 004_seed_payment_data.sql
\i 005_seed_forum_data.sql
\i 006_seed_chat_data.sql

-- Re-enable triggers and foreign key checks
SET session_replication_role = 'origin';

-- Update statistics
ANALYZE;

COMMIT;

-- Display summary of seeded data
SELECT 
    'Users' as table_name, 
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT role::text, ', ') as roles
FROM users
UNION ALL
SELECT 
    'Courses' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status::text, ', ') as status
FROM courses  
UNION ALL
SELECT 
    'Lectures' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status::text, ', ') as status
FROM lectures
UNION ALL
SELECT 
    'Enrollments' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status::text, ', ') as status
FROM enrollments
UNION ALL
SELECT 
    'Payment Methods' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT provider, ', ') as providers
FROM payment_methods
UNION ALL
SELECT 
    'Transactions' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status, ', ') as status
FROM transactions
UNION ALL
SELECT 
    'Subscriptions' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT status, ', ') as status
FROM subscriptions
UNION ALL
SELECT 
    'Forum Topics' as table_name,
    COUNT(*) as record_count,
    STRING_AGG(DISTINCT category, ', ') as categories
FROM forum_topics
UNION ALL
SELECT 
    'Forum Posts' as table_name,
    COUNT(*) as record_count,
    NULL as categories
FROM forum_posts
UNION ALL
SELECT 
    'Chat Sessions' as table_name,
    COUNT(*) as record_count,
    COUNT(CASE WHEN is_active THEN 1 END)::text || ' active' as status
FROM chat_sessions;

-- Display some sample data for verification
SELECT 'Sample Users:' as info;
SELECT username, email, role FROM users LIMIT 5;

SELECT 'Sample Courses:' as info;
SELECT title, instructor_name, category, price FROM courses WHERE status = 'published' LIMIT 3;

SELECT 'Sample Enrollments:' as info;
SELECT u.username, c.title, e.status, e.progress_percentage 
FROM enrollments e 
JOIN users u ON e.user_id = u.id 
JOIN courses c ON e.course_id = c.id 
LIMIT 5;
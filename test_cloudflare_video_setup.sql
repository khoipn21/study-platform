-- Test setup for Cloudflare video streaming with enrollment access control
-- This script sets up test data to validate enrolled course access with video streaming

-- First, let's insert the Cloudflare video into the videos table
INSERT INTO videos (
    id,
    cloudflare_uid,
    title,
    description,
    upload_user_id,
    course_id,
    status,
    visibility,
    stream_url,
    thumbnail_url,
    duration_seconds,
    created_at,
    updated_at
) VALUES (
    'cf111111-1111-1111-1111-111111111111',
    '488bd0b4002f44395a3001d3121dc4f0',
    'Advanced JavaScript Closures - Demo Video',
    'In-depth explanation of JavaScript closures with practical examples',
    '33333333-3333-3333-3333-333333333333', -- Sarah Instructor
    'c8a882cc-6345-4f0d-8562-6e87dc2910ba', -- Advanced JavaScript Concepts course
    'ready',
    'course_only',
    'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/manifest/video.m3u8',
    'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/thumbnails/thumbnail.jpg',
    1800, -- 30 minutes
    NOW(),
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    cloudflare_uid = EXCLUDED.cloudflare_uid,
    title = EXCLUDED.title,
    stream_url = EXCLUDED.stream_url,
    updated_at = NOW();

-- Update the existing lecture to link to our Cloudflare video
UPDATE lectures
SET
    video_id = '488bd0b4002f44395a3001d3121dc4f0',
    video_url = 'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/manifest/video.m3u8',
    duration_minutes = 30
WHERE id = '33333333-3333-3333-3333-333333333001' -- Understanding Closures lecture
AND course_id = 'c8a882cc-6345-4f0d-8562-6e87dc2910ba';

-- Create video quality variants for adaptive streaming
INSERT INTO video_qualities (
    id,
    video_id,
    quality_label,
    bitrate_kbps,
    width,
    height,
    fps,
    codec,
    url,
    created_at
) VALUES
(
    gen_random_uuid(),
    'cf111111-1111-1111-1111-111111111111',
    '360p',
    800,
    640,
    360,
    30,
    'h264',
    'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/manifest/video.m3u8',
    NOW()
),
(
    gen_random_uuid(),
    'cf111111-1111-1111-1111-111111111111',
    '720p',
    2500,
    1280,
    720,
    30,
    'h264',
    'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/manifest/video.m3u8',
    NOW()
),
(
    gen_random_uuid(),
    'cf111111-1111-1111-1111-111111111111',
    '1080p',
    5000,
    1920,
    1080,
    30,
    'h264',
    'https://cloudflarestream.com/488bd0b4002f44395a3001d3121dc4f0/manifest/video.m3u8',
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Add a test user who is NOT enrolled (for negative testing)
INSERT INTO users (
    id,
    email,
    username,
    password_hash,
    role,
    is_email_verified,
    created_at,
    updated_at
) VALUES (
    '11111111-2222-3333-4444-555555555555',
    'teststudent@example.com',
    'teststudent',
    '$2a$10$hash_placeholder_for_testing',
    'student',
    true,
    NOW(),
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    updated_at = NOW();

-- Verify existing enrollments for testing
-- Alice (user_id: 55555555-5555-5555-5555-555555555555) should have access to Advanced JS course
-- Test student (user_id: 11111111-2222-3333-4444-555555555555) should NOT have access

-- Let's also create a free lecture with video for comparison
INSERT INTO videos (
    id,
    cloudflare_uid,
    title,
    description,
    upload_user_id,
    course_id,
    status,
    visibility,
    stream_url,
    thumbnail_url,
    duration_seconds,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    'free_video_placeholder_uid_unique',
    'Introduction to Programming - Welcome Video',
    'Free welcome video for the programming course',
    '22222222-2222-2222-2222-222222222222', -- John Instructor
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', -- Introduction to Programming (free course)
    'ready',
    'public',
    'https://cloudflarestream.com/free_video_placeholder_uid/manifest/video.m3u8',
    'https://cloudflarestream.com/free_video_placeholder_uid/thumbnails/thumbnail.jpg',
    2700, -- 45 minutes
    NOW(),
    NOW()
) ON CONFLICT (cloudflare_uid) DO UPDATE SET
    title = EXCLUDED.title,
    updated_at = NOW();

-- Update the free lecture to link to the free video (using video_id field instead)
UPDATE lectures
SET
    video_id = 'free_video_placeholder_uid_unique',
    video_url = 'https://cloudflarestream.com/free_video_placeholder_uid/manifest/video.m3u8'
WHERE id = '11111111-1111-1111-1111-111111111001' -- What is Programming? lecture
AND course_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- Display test scenario summary
SELECT
    'TEST SCENARIO SETUP COMPLETE' AS status,
    'ENROLLED USERS' AS test_type,
    u.username,
    u.email,
    c.title AS course_title,
    c.price,
    c.currency,
    e.status AS enrollment_status
FROM users u
JOIN enrollments e ON u.id = e.user_id
JOIN courses c ON e.course_id = c.id
WHERE c.id = 'c8a882cc-6345-4f0d-8562-6e87dc2910ba' -- Advanced JavaScript course
UNION ALL
SELECT
    'TEST SCENARIO SETUP COMPLETE' AS status,
    'NON-ENROLLED USER' AS test_type,
    'teststudent' AS username,
    'teststudent@example.com' AS email,
    'Advanced JavaScript Concepts' AS course_title,
    49.99 AS price,
    'USD' AS currency,
    'cancelled' AS enrollment_status;
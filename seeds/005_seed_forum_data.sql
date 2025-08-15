-- Development seed data for forum tables

-- Check if forum tables exist first, then seed data
-- This assumes forum tables from migration 016_create_forum_tables.up.sql

INSERT INTO forum_topics (id, title, description, category, user_id, username, created_at, updated_at, is_pinned, is_locked, post_count, view_count) VALUES

-- Programming discussions
('topic001', 'Best Practices for Learning Programming', 'Share your tips and experiences on how to effectively learn programming', 'General', '55555555-5555-5555-5555-555555555555', 'alice_student', '2024-01-18 10:00:00', '2024-01-25 15:30:00', true, false, 8, 245),
('topic002', 'JavaScript vs Python: Which to Learn First?', 'Comparing JavaScript and Python for beginners', 'Programming', '66666666-6666-6666-6666-666666666666', 'bob_student', '2024-01-20 14:30:00', '2024-02-10 12:20:00', false, false, 12, 189),
('topic003', 'Common Debugging Techniques', 'Share your favorite debugging methods and tools', 'Programming', '22222222-2222-2222-2222-222222222222', 'john_instructor', '2024-01-25 09:15:00', '2024-02-08 16:45:00', false, false, 6, 156),

-- Web development discussions
('topic004', 'React vs Vue: Framework Comparison', 'Pros and cons of different frontend frameworks', 'Web Development', '77777777-7777-7777-7777-777777777777', 'carol_student', '2024-02-01 11:20:00', '2024-02-15 13:25:00', false, false, 15, 312),
('topic005', 'CSS Grid vs Flexbox: When to Use What?', 'Layout techniques in modern CSS', 'Web Development', '33333333-3333-3333-3333-333333333333', 'sarah_instructor', '2024-02-03 16:40:00', '2024-02-20 09:10:00', false, false, 9, 203),

-- Database discussions
('topic006', 'SQL Query Optimization Tips', 'Share your best practices for writing efficient SQL queries', 'Database', '44444444-4444-4444-4444-444444444444', 'mike_instructor', '2024-02-05 08:30:00', '2024-02-18 14:55:00', false, false, 7, 167),

-- Career and general discussions
('topic007', 'From Beginner to Job Ready: Timeline and Tips', 'How long does it take to become job-ready as a developer?', 'Career', '88888888-8888-8888-8888-888888888888', 'david_student', '2024-02-08 13:45:00', '2024-02-22 11:30:00', false, false, 18, 428),
('topic008', 'Course Feedback and Suggestions', 'Share your thoughts on our courses and suggest improvements', 'General', '99999999-9999-9999-9999-999999999999', 'eve_student', '2024-02-12 12:15:00', '2024-02-24 17:20:00', true, false, 11, 289);

-- Seed forum posts
INSERT INTO forum_posts (id, topic_id, user_id, username, content, created_at, updated_at, vote_score, is_solution) VALUES

-- Posts for topic001 (Best Practices for Learning Programming)
('post001', 'topic001', '55555555-5555-5555-5555-555555555555', 'alice_student', 'I''ve found that the best way to learn programming is through hands-on practice. Don''t just watch tutorials - actually code along and experiment with the examples!', '2024-01-18 10:00:00', '2024-01-18 10:00:00', 12, false),
('post002', 'topic001', '66666666-6666-6666-6666-666666666666', 'bob_student', 'Totally agree with Alice! Also, I recommend building small projects to apply what you learn. It helps solidify the concepts.', '2024-01-18 14:20:00', '2024-01-18 14:20:00', 8, false),
('post003', 'topic001', '22222222-2222-2222-2222-222222222222', 'john_instructor', 'Great points! I''d also add that consistent daily practice, even if it''s just 30 minutes, is more effective than long cramming sessions.', '2024-01-19 09:30:00', '2024-01-19 09:30:00', 15, true),

-- Posts for topic002 (JavaScript vs Python)
('post004', 'topic002', '66666666-6666-6666-6666-666666666666', 'bob_student', 'I''m torn between learning JavaScript first for web development or Python for its beginner-friendly syntax. What do you think?', '2024-01-20 14:30:00', '2024-01-20 14:30:00', 5, false),
('post005', 'topic002', '33333333-3333-3333-3333-333333333333', 'sarah_instructor', 'Both are excellent choices! If you''re interested in web development, start with JavaScript. If you''re more interested in data science or general programming, Python might be better.', '2024-01-20 16:45:00', '2024-01-20 16:45:00', 18, true),
('post006', 'topic002', '77777777-7777-7777-7777-777777777777', 'carol_student', 'I started with Python and found it really helped me understand programming fundamentals. Then moving to JavaScript was much easier.', '2024-01-21 11:15:00', '2024-01-21 11:15:00', 9, false),

-- Posts for topic004 (React vs Vue)
('post007', 'topic004', '77777777-7777-7777-7777-777777777777', 'carol_student', 'I''ve been using React for a while, but I keep hearing good things about Vue. What are the main differences?', '2024-02-01 11:20:00', '2024-02-01 11:20:00', 3, false),
('post008', 'topic004', '99999999-9999-9999-9999-999999999999', 'eve_student', 'Vue has a gentler learning curve and more intuitive template syntax. React has a larger ecosystem and more job opportunities. Both are great!', '2024-02-01 15:30:00', '2024-02-01 15:30:00', 14, false),
('post009', 'topic004', '33333333-3333-3333-3333-333333333333', 'sarah_instructor', 'Great question! React''s component model and hooks are powerful, while Vue''s single-file components are very developer-friendly. Choose based on your team and project needs.', '2024-02-02 10:40:00', '2024-02-02 10:40:00', 22, true),

-- Posts for topic007 (Career timeline)
('post010', 'topic007', '88888888-8888-8888-8888-888888888888', 'david_student', 'How long did it take everyone to land their first developer job? I''ve been learning for 8 months and starting to feel ready to apply.', '2024-02-08 13:45:00', '2024-02-08 13:45:00', 7, false),
('post011', 'topic007', '22222222-2222-2222-2222-222222222222', 'john_instructor', 'The timeline varies greatly, but 6-12 months of consistent learning is typical. Focus on building a strong portfolio with 3-4 solid projects.', '2024-02-08 16:20:00', '2024-02-08 16:20:00', 16, false),
('post012', 'topic007', '55555555-5555-5555-5555-555555555555', 'alice_student', 'It took me about 10 months, but I was studying part-time. The key is not just learning syntax, but understanding how to solve problems and build complete applications.', '2024-02-09 09:10:00', '2024-02-09 09:10:00', 11, true);

-- Seed forum votes
INSERT INTO forum_votes (id, post_id, user_id, vote_type, created_at) VALUES

-- Votes for various posts
('vote001', 'post001', '66666666-6666-6666-6666-666666666666', 'upvote', '2024-01-18 15:00:00'),
('vote002', 'post001', '77777777-7777-7777-7777-777777777777', 'upvote', '2024-01-18 16:30:00'),
('vote003', 'post001', '88888888-8888-8888-8888-888888888888', 'upvote', '2024-01-19 08:45:00'),
('vote004', 'post003', '55555555-5555-5555-5555-555555555555', 'upvote', '2024-01-19 10:15:00'),
('vote005', 'post003', '88888888-8888-8888-8888-888888888888', 'upvote', '2024-01-19 14:20:00'),
('vote006', 'post005', '66666666-6666-6666-6666-666666666666', 'upvote', '2024-01-20 17:00:00'),
('vote007', 'post005', '77777777-7777-7777-7777-777777777777', 'upvote', '2024-01-21 09:30:00'),
('vote008', 'post009', '77777777-7777-7777-7777-777777777777', 'upvote', '2024-02-02 11:00:00'),
('vote009', 'post009', '88888888-8888-8888-8888-888888888888', 'upvote', '2024-02-02 15:45:00'),
('vote010', 'post011', '99999999-9999-9999-9999-999999999999', 'upvote', '2024-02-08 18:30:00');
-- Development seed data for forum tables (corrected for actual schema)

-- Seed forum topics
INSERT INTO forum_topics (id, title, course_id, creator_id, is_pinned, is_locked, created_at, updated_at) VALUES

-- General programming discussions  
('11111111-1111-1111-1111-111111111111', 'Best Practices for Learning Programming', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', true, false, '2024-01-18 10:00:00', '2024-01-25 15:30:00'),
('22222222-2222-2222-2222-222222222222', 'JavaScript vs Python: Which to Learn First?', NULL, '66666666-6666-6666-6666-666666666666', false, false, '2024-01-20 14:30:00', '2024-02-10 12:20:00'),
('33333333-3333-3333-3333-333333333333', 'Common Debugging Techniques', NULL, '22222222-2222-2222-2222-222222222222', false, false, '2024-01-25 09:15:00', '2024-02-08 16:45:00'),

-- Course-specific discussions
('44444444-4444-4444-4444-444444444444', 'React vs Vue: Framework Comparison', 'dddddddd-dddd-dddd-dddd-dddddddddddd', '77777777-7777-7777-7777-777777777777', false, false, '2024-02-01 11:20:00', '2024-02-15 13:25:00'),
('55555555-5555-5555-5555-555555555555', 'CSS Grid vs Flexbox: When to Use What?', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '33333333-3333-3333-3333-333333333333', false, false, '2024-02-03 16:40:00', '2024-02-20 09:10:00'),
('66666666-6666-6666-6666-666666666666', 'SQL Query Optimization Tips', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '44444444-4444-4444-4444-444444444444', false, false, '2024-02-05 08:30:00', '2024-02-18 14:55:00'),

-- General discussions
('77777777-7777-7777-7777-777777777777', 'From Beginner to Job Ready: Timeline and Tips', NULL, '88888888-8888-8888-8888-888888888888', false, false, '2024-02-08 13:45:00', '2024-02-22 11:30:00'),
('88888888-8888-8888-8888-888888888888', 'Course Feedback and Suggestions', NULL, '99999999-9999-9999-9999-999999999999', true, false, '2024-02-12 12:15:00', '2024-02-24 17:20:00');

-- Seed forum posts
INSERT INTO forum_posts (id, topic_id, user_id, content, is_solution, created_at, updated_at) VALUES

-- Posts for topic001 (Best Practices for Learning Programming)
('11111111-1111-1111-1111-111111111001', '11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555', 'I''ve found that the best way to learn programming is through hands-on practice. Don''t just watch tutorials - actually code along and experiment with the examples!', false, '2024-01-18 10:00:00', '2024-01-18 10:00:00'),
('11111111-1111-1111-1111-111111111002', '11111111-1111-1111-1111-111111111111', '66666666-6666-6666-6666-666666666666', 'Totally agree with Alice! Also, I recommend building small projects to apply what you learn. It helps solidify the concepts.', false, '2024-01-18 14:20:00', '2024-01-18 14:20:00'),
('11111111-1111-1111-1111-111111111003', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Great points! I''d also add that consistent daily practice, even if it''s just 30 minutes, is more effective than long cramming sessions.', true, '2024-01-19 09:30:00', '2024-01-19 09:30:00'),

-- Posts for topic002 (JavaScript vs Python)
('22222222-2222-2222-2222-222222222001', '22222222-2222-2222-2222-222222222222', '66666666-6666-6666-6666-666666666666', 'I''m torn between learning JavaScript first for web development or Python for its beginner-friendly syntax. What do you think?', false, '2024-01-20 14:30:00', '2024-01-20 14:30:00'),
('22222222-2222-2222-2222-222222222002', '22222222-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333', 'Both are excellent choices! If you''re interested in web development, start with JavaScript. If you''re more interested in data science or general programming, Python might be better.', true, '2024-01-20 16:45:00', '2024-01-20 16:45:00'),
('22222222-2222-2222-2222-222222222003', '22222222-2222-2222-2222-222222222222', '77777777-7777-7777-7777-777777777777', 'I started with Python and found it really helped me understand programming fundamentals. Then moving to JavaScript was much easier.', false, '2024-01-21 11:15:00', '2024-01-21 11:15:00'),

-- Posts for topic004 (React vs Vue)
('44444444-4444-4444-4444-444444444001', '44444444-4444-4444-4444-444444444444', '77777777-7777-7777-7777-777777777777', 'I''ve been using React for a while, but I keep hearing good things about Vue. What are the main differences?', false, '2024-02-01 11:20:00', '2024-02-01 11:20:00'),
('44444444-4444-4444-4444-444444444002', '44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999999', 'Vue has a gentler learning curve and more intuitive template syntax. React has a larger ecosystem and more job opportunities. Both are great!', false, '2024-02-01 15:30:00', '2024-02-01 15:30:00'),
('44444444-4444-4444-4444-444444444003', '44444444-4444-4444-4444-444444444444', '33333333-3333-3333-3333-333333333333', 'Great question! React''s component model and hooks are powerful, while Vue''s single-file components are very developer-friendly. Choose based on your team and project needs.', true, '2024-02-02 10:40:00', '2024-02-02 10:40:00'),

-- Posts for topic007 (Career timeline)
('77777777-7777-7777-7777-777777777001', '77777777-7777-7777-7777-777777777777', '88888888-8888-8888-8888-888888888888', 'How long did it take everyone to land their first developer job? I''ve been learning for 8 months and starting to feel ready to apply.', false, '2024-02-08 13:45:00', '2024-02-08 13:45:00'),
('77777777-7777-7777-7777-777777777002', '77777777-7777-7777-7777-777777777777', '22222222-2222-2222-2222-222222222222', 'The timeline varies greatly, but 6-12 months of consistent learning is typical. Focus on building a strong portfolio with 3-4 solid projects.', false, '2024-02-08 16:20:00', '2024-02-08 16:20:00'),
('77777777-7777-7777-7777-777777777003', '77777777-7777-7777-7777-777777777777', '55555555-5555-5555-5555-555555555555', 'It took me about 10 months, but I was studying part-time. The key is not just learning syntax, but understanding how to solve problems and build complete applications.', true, '2024-02-09 09:10:00', '2024-02-09 09:10:00');
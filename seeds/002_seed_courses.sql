-- Development seed data for courses and lectures

-- Seed courses
INSERT INTO courses (id, title, description, instructor_id, instructor_name, category, level, price, currency, thumbnail_url, status, duration_minutes, enrollment_count, rating, rating_count, tags, created_at, updated_at) VALUES

-- Free courses
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Introduction to Programming', 'Learn the fundamentals of programming with hands-on examples and exercises. Perfect for beginners who want to start their coding journey.', '22222222-2222-2222-2222-222222222222', 'John Instructor', 'Programming', 'beginner', 0, 'USD', 'https://images.unsplash.com/photo-1516116216624-53e697fedbea', 'published', 480, 150, 4.5, 32, '{"programming", "beginner", "fundamentals", "coding"}', '2024-01-15 10:00:00', '2024-01-15 10:00:00'),

('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Web Development Basics', 'Build your first website using HTML, CSS, and JavaScript. Learn responsive design and modern web development practices.', '22222222-2222-2222-2222-222222222222', 'John Instructor', 'Web Development', 'beginner', 0, 'USD', 'https://images.unsplash.com/photo-1547658719-da2b51169166', 'published', 600, 89, 4.3, 18, '{"web", "html", "css", "javascript", "responsive"}', '2024-01-20 11:00:00', '2024-01-20 11:00:00'),

-- Paid courses
('c8a882cc-6345-4f0d-8562-6e87dc2910ba', 'Advanced JavaScript Concepts', 'Deep dive into advanced JavaScript topics including closures, prototypes, async/await, and modern ES6+ features.', '33333333-3333-3333-3333-333333333333', 'Sarah Instructor', 'Programming', 'advanced', 49.99, 'USD', 'https://images.unsplash.com/photo-1579468118864-1b9ea3c0db4a', 'published', 720, 67, 4.8, 25, '{"javascript", "advanced", "es6", "async", "closures"}', '2024-02-01 09:00:00', '2024-02-01 09:00:00'),

('dddddddd-dddd-dddd-dddd-dddddddddddd', 'React Development Masterclass', 'Master React.js from basics to advanced concepts. Build real-world applications with hooks, context, and modern patterns.', '33333333-3333-3333-3333-333333333333', 'Sarah Instructor', 'Web Development', 'intermediate', 79.99, 'USD', 'https://images.unsplash.com/photo-1633356122544-f134324a6cee', 'published', 900, 45, 4.7, 15, '{"react", "javascript", "hooks", "components", "spa"}', '2024-02-10 14:00:00', '2024-02-10 14:00:00'),

('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Database Design & SQL', 'Learn database design principles and master SQL queries. Cover PostgreSQL, indexes, optimization, and best practices.', '44444444-4444-4444-4444-444444444444', 'Mike Instructor', 'Database', 'intermediate', 39.99, 'USD', 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31', 'published', 540, 78, 4.4, 22, '{"sql", "database", "postgresql", "design", "optimization"}', '2024-02-15 13:00:00', '2024-02-15 13:00:00'),

('ffffffff-ffff-ffff-ffff-ffffffffffff', 'Mobile App Development with React Native', 'Build cross-platform mobile applications using React Native. Learn navigation, state management, and deployment.', '44444444-4444-4444-4444-444444444444', 'Mike Instructor', 'Mobile Development', 'intermediate', 89.99, 'USD', 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c', 'published', 1080, 34, 4.6, 12, '{"react-native", "mobile", "ios", "android", "cross-platform"}', '2024-02-20 16:00:00', '2024-02-20 16:00:00'),

-- Draft courses
('10101010-1010-1010-1010-101010101010', 'Machine Learning Fundamentals', 'Introduction to machine learning algorithms and practical applications using Python and popular ML libraries.', '22222222-2222-2222-2222-222222222222', 'John Instructor', 'Data Science', 'intermediate', 99.99, 'USD', 'https://images.unsplash.com/photo-1551288049-bebda4e38f71', 'draft', 0, 0, 0, 0, '{"machine-learning", "python", "ai", "data-science"}', '2024-02-25 12:00:00', '2024-02-25 12:00:00');

-- Seed lectures for courses
INSERT INTO lectures (id, course_id, title, description, order_number, duration_minutes, video_url, video_id, status, is_free, created_at, updated_at) VALUES

-- Lectures for Introduction to Programming
('11111111-1111-1111-1111-111111111001', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'What is Programming?', 'Understanding the basics of programming and why it matters', 1, 45, 'https://example.com/video1', 'vid_001', 'published', true, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),
('11111111-1111-1111-1111-111111111002', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Setting Up Your Development Environment', 'Installing and configuring your first code editor and tools', 2, 30, 'https://example.com/video2', 'vid_002', 'published', true, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),
('11111111-1111-1111-1111-111111111003', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Variables and Data Types', 'Learning about different types of data and how to store them', 3, 40, 'https://example.com/video3', 'vid_003', 'published', false, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),
('11111111-1111-1111-1111-111111111004', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Control Flow: If Statements', 'Making decisions in your code with conditional statements', 4, 35, 'https://example.com/video4', 'vid_004', 'published', false, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),

-- Lectures for Web Development Basics
('22222222-2222-2222-2222-222222222001', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'HTML Fundamentals', 'Building the structure of web pages with HTML', 1, 50, 'https://example.com/video5', 'vid_005', 'published', true, '2024-01-20 11:30:00', '2024-01-20 11:30:00'),
('22222222-2222-2222-2222-222222222002', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'CSS Styling and Layout', 'Making your websites look beautiful with CSS', 2, 60, 'https://example.com/video6', 'vid_006', 'published', false, '2024-01-20 11:30:00', '2024-01-20 11:30:00'),
('22222222-2222-2222-2222-222222222003', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'JavaScript Interactivity', 'Adding dynamic behavior to your web pages', 3, 55, 'https://example.com/video7', 'vid_007', 'published', false, '2024-01-20 11:30:00', '2024-01-20 11:30:00'),

-- Lectures for Advanced JavaScript Concepts  
('33333333-3333-3333-3333-333333333001', 'c8a882cc-6345-4f0d-8562-6e87dc2910ba', 'Understanding Closures', 'Deep dive into JavaScript closures and lexical scoping', 1, 45, 'https://example.com/video8', 'vid_008', 'published', true, '2024-02-01 09:30:00', '2024-02-01 09:30:00'),
('33333333-3333-3333-3333-333333333002', 'c8a882cc-6345-4f0d-8562-6e87dc2910ba', 'Prototypes and Inheritance', 'Mastering JavaScript object-oriented programming', 2, 50, 'https://example.com/video9', 'vid_009', 'published', false, '2024-02-01 09:30:00', '2024-02-01 09:30:00'),
('33333333-3333-3333-3333-333333333003', 'c8a882cc-6345-4f0d-8562-6e87dc2910ba', 'Async/Await and Promises', 'Handling asynchronous operations in modern JavaScript', 3, 60, 'https://example.com/video10', 'vid_010', 'published', false, '2024-02-01 09:30:00', '2024-02-01 09:30:00'),

-- Lectures for React Development Masterclass
('44444444-4444-4444-4444-444444444001', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'React Components and JSX', 'Building your first React components', 1, 40, 'https://example.com/video11', 'vid_011', 'published', true, '2024-02-10 14:30:00', '2024-02-10 14:30:00'),
('44444444-4444-4444-4444-444444444002', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'State and Props', 'Managing data flow in React applications', 2, 45, 'https://example.com/video12', 'vid_012', 'published', false, '2024-02-10 14:30:00', '2024-02-10 14:30:00'),
('44444444-4444-4444-4444-444444444003', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'React Hooks Deep Dive', 'useState, useEffect, and custom hooks', 3, 55, 'https://example.com/video13', 'vid_013', 'published', false, '2024-02-10 14:30:00', '2024-02-10 14:30:00'),

-- Lectures for Database Design & SQL
('55555555-5555-5555-5555-555555555001', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Database Design Principles', 'Normalization, relationships, and best practices', 1, 50, 'https://example.com/video14', 'vid_014', 'published', true, '2024-02-15 13:30:00', '2024-02-15 13:30:00'),
('55555555-5555-5555-5555-555555555002', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'SQL Fundamentals', 'SELECT, INSERT, UPDATE, DELETE operations', 2, 40, 'https://example.com/video15', 'vid_015', 'published', false, '2024-02-15 13:30:00', '2024-02-15 13:30:00'),
('55555555-5555-5555-5555-555555555003', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Advanced SQL Queries', 'JOINs, subqueries, and complex operations', 3, 60, 'https://example.com/video16', 'vid_016', 'published', false, '2024-02-15 13:30:00', '2024-02-15 13:30:00');
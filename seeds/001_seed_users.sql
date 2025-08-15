-- Development seed data for users table
-- Note: All passwords are hashed version of 'password123'

INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at) VALUES
-- Admin users
('11111111-1111-1111-1111-111111111111', 'admin', 'admin@studyplatform.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),

-- Instructor users  
('22222222-2222-2222-2222-222222222222', 'john_instructor', 'john@studyplatform.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'instructor', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('33333333-3333-3333-3333-333333333333', 'sarah_instructor', 'sarah@studyplatform.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'instructor', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('44444444-4444-4444-4444-444444444444', 'mike_instructor', 'mike@studyplatform.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'instructor', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),

-- Student users
('55555555-5555-5555-5555-555555555555', 'alice_student', 'alice@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('66666666-6666-6666-6666-666666666666', 'bob_student', 'bob@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('77777777-7777-7777-7777-777777777777', 'carol_student', 'carol@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('88888888-8888-8888-8888-888888888888', 'david_student', 'david@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('99999999-9999-9999-9999-999999999999', 'eve_student', 'eve@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- Test user for integration testing
INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at) VALUES
('fc7f125f-f322-458c-b789-5a209ffa1fed', 'test_user', 'test@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', '2024-01-01 00:00:00', '2024-01-01 00:00:00');
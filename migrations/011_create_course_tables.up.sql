-- Course status enum
CREATE TYPE course_status AS ENUM ('draft', 'published', 'archived');

-- Course level enum
CREATE TYPE course_level AS ENUM ('beginner', 'intermediate', 'advanced');

-- Lecture status enum
CREATE TYPE lecture_status AS ENUM ('draft', 'published');

-- Enrollment status enum
CREATE TYPE enrollment_status AS ENUM ('enrolled', 'completed', 'cancelled');

-- Create courses table
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    instructor_id UUID NOT NULL,
    instructor_name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    level course_level NOT NULL DEFAULT 'beginner',
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    thumbnail_url TEXT,
    status course_status NOT NULL DEFAULT 'draft',
    duration_minutes INTEGER NOT NULL DEFAULT 0,
    enrollment_count INTEGER NOT NULL DEFAULT 0,
    rating DECIMAL(3,2) NOT NULL DEFAULT 0,
    rating_count INTEGER NOT NULL DEFAULT 0,
    tags TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create lectures table
CREATE TABLE lectures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_number INTEGER NOT NULL,
    duration_minutes INTEGER NOT NULL DEFAULT 0,
    video_url TEXT,
    video_id VARCHAR(255),
    status lecture_status NOT NULL DEFAULT 'draft',
    is_free BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create enrollments table
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status enrollment_status NOT NULL DEFAULT 'enrolled',
    progress_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,
    enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    last_accessed TIMESTAMP,
    UNIQUE(user_id, course_id)
);

-- Create indexes for courses
CREATE INDEX idx_courses_instructor_id ON courses(instructor_id);
CREATE INDEX idx_courses_category ON courses(category);
CREATE INDEX idx_courses_level ON courses(level);
CREATE INDEX idx_courses_status ON courses(status);
CREATE INDEX idx_courses_price ON courses(price);
CREATE INDEX idx_courses_rating ON courses(rating);
CREATE INDEX idx_courses_created_at ON courses(created_at);
CREATE INDEX idx_courses_tags ON courses USING gin(tags);

-- Create indexes for lectures
CREATE INDEX idx_lectures_course_id ON lectures(course_id);
CREATE INDEX idx_lectures_status ON lectures(status);
CREATE INDEX idx_lectures_order_number ON lectures(course_id, order_number);

-- Create indexes for enrollments
CREATE INDEX idx_enrollments_user_id ON enrollments(user_id);
CREATE INDEX idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX idx_enrollments_status ON enrollments(status);
CREATE INDEX idx_enrollments_enrolled_at ON enrollments(enrolled_at);
CREATE INDEX idx_enrollments_progress ON enrollments(progress_percentage);

-- Create full-text search index for courses
CREATE INDEX idx_courses_search ON courses USING gin(to_tsvector('english', title || ' ' || description || ' ' || category));

-- Add foreign key constraint for instructor_id (references users table)
ALTER TABLE courses ADD CONSTRAINT fk_courses_instructor_id FOREIGN KEY (instructor_id) REFERENCES users(id);

-- Add foreign key constraint for enrollment user_id (references users table)
ALTER TABLE enrollments ADD CONSTRAINT fk_enrollments_user_id FOREIGN KEY (user_id) REFERENCES users(id);
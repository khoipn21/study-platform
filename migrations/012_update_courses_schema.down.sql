-- Drop the updated tables
DROP TABLE IF EXISTS enrollments CASCADE;
DROP TABLE IF EXISTS lectures CASCADE;
DROP TABLE IF EXISTS courses CASCADE;

-- Drop custom types
DROP TYPE IF EXISTS enrollment_status;
DROP TYPE IF EXISTS lecture_status;
DROP TYPE IF EXISTS course_level;
DROP TYPE IF EXISTS course_status;

-- Recreate original courses table
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    creator_id UUID REFERENCES users(id),
    thumbnail_url TEXT,
    price DECIMAL(10,2) DEFAULT 0.00,
    is_free BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_creator_id ON courses(creator_id);
CREATE INDEX idx_courses_is_free ON courses(is_free);
CREATE INDEX idx_courses_created_at ON courses(created_at);

-- Recreate original enrollment table
CREATE TABLE enrollment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    progress_percentage DECIMAL(5,2) DEFAULT 0.00,
    UNIQUE(user_id, course_id)
);

CREATE INDEX idx_enrollment_user_id ON enrollment(user_id);
CREATE INDEX idx_enrollment_course_id ON enrollment(course_id);
CREATE INDEX idx_enrollment_enrolled_at ON enrollment(enrolled_at);
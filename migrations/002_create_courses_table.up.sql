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
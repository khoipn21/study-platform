CREATE TABLE lectures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    video_url TEXT,
    duration INT,
    sequence_order INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lectures_course_id ON lectures(course_id);
CREATE INDEX idx_lectures_sequence_order ON lectures(sequence_order);
CREATE UNIQUE INDEX idx_lectures_course_sequence ON lectures(course_id, sequence_order);
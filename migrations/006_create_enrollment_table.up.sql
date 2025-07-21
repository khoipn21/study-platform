CREATE TABLE enrollment (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    transaction_id UUID,
    enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (user_id, course_id)
);

CREATE INDEX idx_enrollment_user_id ON enrollment(user_id);
CREATE INDEX idx_enrollment_course_id ON enrollment(course_id);
CREATE INDEX idx_enrollment_is_active ON enrollment(is_active);
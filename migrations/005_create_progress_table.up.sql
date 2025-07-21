CREATE TABLE progress (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    lecture_id UUID REFERENCES lectures(id) ON DELETE CASCADE,
    watched_duration INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    last_watched_at TIMESTAMP,
    PRIMARY KEY (user_id, lecture_id)
);

CREATE INDEX idx_progress_user_id ON progress(user_id);
CREATE INDEX idx_progress_lecture_id ON progress(lecture_id);
CREATE INDEX idx_progress_completed ON progress(completed);
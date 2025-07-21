-- Update progress table to match progress service requirements
ALTER TABLE progress 
ADD COLUMN IF NOT EXISTS is_completed BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS last_accessed TIMESTAMP DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_progress_user_course ON progress(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_progress_lecture_id ON progress(lecture_id);
CREATE INDEX IF NOT EXISTS idx_progress_is_completed ON progress(is_completed);
CREATE INDEX IF NOT EXISTS idx_progress_last_accessed ON progress(last_accessed);
CREATE INDEX IF NOT EXISTS idx_progress_completed_at ON progress(completed_at);

-- Update enrollments table structure to match progress service
ALTER TABLE enrollments 
ADD COLUMN IF NOT EXISTS completed_lectures INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_lectures INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_watch_time_seconds INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS last_accessed TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Create indexes for enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_created_at ON enrollments(created_at);
CREATE INDEX IF NOT EXISTS idx_enrollments_updated_at ON enrollments(updated_at);
CREATE INDEX IF NOT EXISTS idx_enrollments_last_accessed ON enrollments(last_accessed);
CREATE INDEX IF NOT EXISTS idx_enrollments_completed_lectures ON enrollments(completed_lectures);
CREATE INDEX IF NOT EXISTS idx_enrollments_total_lectures ON enrollments(total_lectures);

-- Create a trigger to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply the trigger to progress table
DROP TRIGGER IF EXISTS update_progress_updated_at ON progress;
CREATE TRIGGER update_progress_updated_at
    BEFORE UPDATE ON progress
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply the trigger to enrollments table
DROP TRIGGER IF EXISTS update_enrollments_updated_at ON enrollments;
CREATE TRIGGER update_enrollments_updated_at
    BEFORE UPDATE ON enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
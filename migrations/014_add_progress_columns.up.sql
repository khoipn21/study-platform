-- Add missing columns to progress table
ALTER TABLE progress 
ADD COLUMN IF NOT EXISTS id UUID DEFAULT gen_random_uuid(),
ADD COLUMN IF NOT EXISTS course_id UUID,
ADD COLUMN IF NOT EXISTS progress_percentage DECIMAL(5,2) DEFAULT 0.0,
ADD COLUMN IF NOT EXISTS watch_time_seconds INTEGER DEFAULT 0;

-- Update course_id from lecture_id by joining with lectures table
UPDATE progress 
SET course_id = l.course_id
FROM lectures l
WHERE progress.lecture_id = l.id AND progress.course_id IS NULL;

-- Make course_id NOT NULL after populating it
ALTER TABLE progress 
ALTER COLUMN course_id SET NOT NULL;

-- Add foreign key constraint for course_id
ALTER TABLE progress 
ADD CONSTRAINT fk_progress_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;

-- Create composite index for better performance
CREATE INDEX IF NOT EXISTS idx_progress_user_course ON progress(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_progress_course_id ON progress(course_id);
CREATE INDEX IF NOT EXISTS idx_progress_progress_percentage ON progress(progress_percentage);
CREATE INDEX IF NOT EXISTS idx_progress_watch_time ON progress(watch_time_seconds);

-- Update progress_percentage based on existing completed status
UPDATE progress 
SET progress_percentage = CASE 
    WHEN completed = true OR is_completed = true THEN 100.0
    ELSE 0.0
END
WHERE progress_percentage = 0.0;

-- Update watch_time_seconds based on existing watched_duration
UPDATE progress 
SET watch_time_seconds = COALESCE(watched_duration, 0)
WHERE watch_time_seconds = 0;
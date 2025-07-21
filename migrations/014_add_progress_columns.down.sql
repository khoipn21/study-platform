-- Remove new columns from progress table
ALTER TABLE progress 
DROP CONSTRAINT IF EXISTS fk_progress_course_id,
DROP COLUMN IF EXISTS progress_percentage,
DROP COLUMN IF EXISTS watch_time_seconds,
DROP COLUMN IF EXISTS course_id,
DROP COLUMN IF EXISTS id;

-- Remove indexes
DROP INDEX IF EXISTS idx_progress_user_course;
DROP INDEX IF EXISTS idx_progress_course_id;
DROP INDEX IF EXISTS idx_progress_progress_percentage;
DROP INDEX IF EXISTS idx_progress_watch_time;
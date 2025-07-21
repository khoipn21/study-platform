-- Drop triggers
DROP TRIGGER IF EXISTS update_enrollments_updated_at ON enrollments;
DROP TRIGGER IF EXISTS update_progress_updated_at ON progress;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Remove columns from enrollments table
ALTER TABLE enrollments 
DROP COLUMN IF EXISTS completed_lectures,
DROP COLUMN IF EXISTS total_lectures,
DROP COLUMN IF EXISTS total_watch_time_seconds,
DROP COLUMN IF EXISTS last_accessed,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

-- Remove indexes from enrollments
DROP INDEX IF EXISTS idx_enrollments_created_at;
DROP INDEX IF EXISTS idx_enrollments_updated_at;
DROP INDEX IF EXISTS idx_enrollments_last_accessed;
DROP INDEX IF EXISTS idx_enrollments_completed_lectures;
DROP INDEX IF EXISTS idx_enrollments_total_lectures;

-- Remove columns from progress table
ALTER TABLE progress 
DROP COLUMN IF EXISTS is_completed,
DROP COLUMN IF EXISTS last_accessed,
DROP COLUMN IF EXISTS completed_at,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

-- Remove indexes from progress table
DROP INDEX IF EXISTS idx_progress_user_course;
DROP INDEX IF EXISTS idx_progress_lecture_id;
DROP INDEX IF EXISTS idx_progress_is_completed;
DROP INDEX IF EXISTS idx_progress_last_accessed;
DROP INDEX IF EXISTS idx_progress_completed_at;
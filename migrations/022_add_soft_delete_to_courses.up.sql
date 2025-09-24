-- Add soft delete support to courses table
ALTER TABLE courses ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Create index for soft delete queries (only include non-deleted records)
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON courses(deleted_at) WHERE deleted_at IS NULL;

-- Add soft delete support to lectures table as well (cascade soft delete)
ALTER TABLE lectures ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Create index for soft delete queries on lectures
CREATE INDEX IF NOT EXISTS idx_lectures_deleted_at ON lectures(deleted_at) WHERE deleted_at IS NULL;

-- Add soft delete support to enrollments table (for data integrity)
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Create index for soft delete queries on enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_deleted_at ON enrollments(deleted_at) WHERE deleted_at IS NULL;

-- Update existing queries to exclude soft-deleted records by default
-- This is handled in the application layer, but we add comments for clarity

-- Note: Application code will need to be updated to:
-- 1. Add WHERE deleted_at IS NULL to all SELECT queries
-- 2. Use UPDATE courses SET deleted_at = NOW() instead of DELETE
-- 3. Cascade soft delete to related lectures and enrollments
-- Remove soft delete support from courses table
DROP INDEX IF EXISTS idx_courses_deleted_at;
ALTER TABLE courses DROP COLUMN IF EXISTS deleted_at;

-- Remove soft delete support from lectures table
DROP INDEX IF EXISTS idx_lectures_deleted_at;
ALTER TABLE lectures DROP COLUMN IF EXISTS deleted_at;

-- Remove soft delete support from enrollments table
DROP INDEX IF EXISTS idx_enrollments_deleted_at;
ALTER TABLE enrollments DROP COLUMN IF EXISTS deleted_at;
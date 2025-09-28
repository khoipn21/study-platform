-- =============================================
-- ROLLBACK: MIGRATE RESOURCES FROM COURSES TO LECTURES
-- Migration: 025_migrate_resources_to_lectures (DOWN)
-- =============================================

-- Drop the resource distribution view
DROP VIEW IF EXISTS resource_distribution_view;

-- Remove migration log entry
DELETE FROM migration_logs WHERE migration_name = '025_migrate_resources_to_lectures';

-- Remove the additional constraints and columns from course_resource_analytics
ALTER TABLE course_resource_analytics DROP CONSTRAINT IF EXISTS check_analytics_reference;
ALTER TABLE course_resource_analytics DROP CONSTRAINT IF EXISTS uk_course_resource_analytics_unique;

-- Restore original unique constraint for course_resource_analytics
ALTER TABLE course_resource_analytics
ADD CONSTRAINT course_resource_analytics_course_id_resource_id_analytics_date_key
    UNIQUE (course_id, resource_id, analytics_date);

-- Drop the lecture analytics index
DROP INDEX IF EXISTS idx_course_resource_analytics_lecture_date;

-- Remove lecture_id column from course_resource_analytics
ALTER TABLE course_resource_analytics DROP COLUMN IF EXISTS lecture_id;

-- Remove comment updates
COMMENT ON TABLE course_resource_analytics IS 'Analytics for course resources like videos, documents, etc.';

-- Remove trigger and function that prevent course_resources insertions
DROP TRIGGER IF EXISTS trigger_prevent_course_resources_insert ON course_resources;
DROP FUNCTION IF EXISTS prevent_course_resources_insert();

-- Remove deprecation comment from course_resources table
COMMENT ON TABLE course_resources IS NULL;

-- Drop lecture_resources table and all its constraints/indexes
DROP INDEX IF EXISTS idx_lecture_resources_type;
DROP INDEX IF EXISTS idx_lecture_resources_order;
DROP INDEX IF EXISTS idx_lecture_resources_file_id;
DROP INDEX IF EXISTS idx_lecture_resources_lecture_id;

-- Drop the table (this will automatically drop foreign key constraints)
DROP TABLE IF EXISTS lecture_resources;
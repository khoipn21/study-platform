-- =============================================
-- MIGRATE RESOURCES FROM COURSES TO LECTURES
-- Migration: 025_migrate_resources_to_lectures
-- =============================================

-- Create lecture_resources table
CREATE TABLE lecture_resources (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecture_id    UUID NOT NULL,
    file_id       UUID NOT NULL,
    resource_type VARCHAR(50) NOT NULL DEFAULT 'document',
    display_order INTEGER NOT NULL DEFAULT 1,
    is_required   BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Foreign key constraints
    CONSTRAINT fk_lecture_resources_lecture_id
        FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
    CONSTRAINT fk_lecture_resources_file_id
        FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,

    -- Unique constraint to prevent duplicate resource assignments
    CONSTRAINT uk_lecture_resources_lecture_file
        UNIQUE (lecture_id, file_id)
);

-- Create indexes for performance
CREATE INDEX idx_lecture_resources_lecture_id ON lecture_resources(lecture_id);
CREATE INDEX idx_lecture_resources_file_id ON lecture_resources(file_id);
CREATE INDEX idx_lecture_resources_order ON lecture_resources(lecture_id, display_order);
CREATE INDEX idx_lecture_resources_type ON lecture_resources(resource_type);

-- Add resource type validation
ALTER TABLE lecture_resources ADD CONSTRAINT check_resource_type
    CHECK (resource_type IN (
        'document', 'pdf', 'video', 'audio', 'image', 'archive',
        'code', 'slides', 'worksheet', 'quiz', 'assignment', 'other'
    ));

-- Add comment to table for documentation
COMMENT ON TABLE lecture_resources IS 'Resources attached to individual lectures instead of courses';
COMMENT ON COLUMN lecture_resources.lecture_id IS 'Reference to the lecture this resource belongs to';
COMMENT ON COLUMN lecture_resources.file_id IS 'Reference to the file in bucket service';
COMMENT ON COLUMN lecture_resources.resource_type IS 'Type of resource (document, pdf, video, etc.)';
COMMENT ON COLUMN lecture_resources.display_order IS 'Order in which resources should be displayed';
COMMENT ON COLUMN lecture_resources.is_required IS 'Whether this resource is required for course completion';

-- Migration Strategy: Since course_resources is empty, we don't need to migrate data
-- But we'll add a trigger to prevent accidental data insertion into the old table

-- Create trigger function to prevent new course_resources
CREATE OR REPLACE FUNCTION prevent_course_resources_insert()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Resources should now be attached to lectures, not courses. Use lecture_resources table instead.';
END;
$$ LANGUAGE plpgsql;

-- Create trigger on course_resources to prevent new insertions
CREATE TRIGGER trigger_prevent_course_resources_insert
    BEFORE INSERT ON course_resources
    FOR EACH ROW
    EXECUTE FUNCTION prevent_course_resources_insert();

-- Add a soft-deprecation comment to the old table
COMMENT ON TABLE course_resources IS 'DEPRECATED: Use lecture_resources instead. Resources are now attached to individual lectures.';

-- Update course_resource_analytics table to support both course and lecture level analytics (if it exists)
-- Add lecture_id column to track lecture-specific resource analytics
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'course_resource_analytics') THEN
        ALTER TABLE course_resource_analytics
        ADD COLUMN lecture_id UUID REFERENCES lectures(id) ON DELETE CASCADE;

        -- Create index for lecture-based analytics queries
        CREATE INDEX IF NOT EXISTS idx_course_resource_analytics_lecture_date
            ON course_resource_analytics(lecture_id, analytics_date);

        -- Update the constraint to allow both course-level and lecture-level analytics
        ALTER TABLE course_resource_analytics
        DROP CONSTRAINT IF EXISTS course_resource_analytics_course_id_resource_id_analytics_date_key;

        -- Add new unique constraint that handles both cases
        ALTER TABLE course_resource_analytics
        ADD CONSTRAINT uk_course_resource_analytics_unique
            UNIQUE (course_id, resource_id, lecture_id, analytics_date);

        -- Add check constraint to ensure either course_id or lecture_id is set (or both)
        ALTER TABLE course_resource_analytics
        ADD CONSTRAINT check_analytics_reference
            CHECK (course_id IS NOT NULL OR lecture_id IS NOT NULL);

        -- Update comment for analytics table
        COMMENT ON TABLE course_resource_analytics IS 'Analytics for course and lecture resources. Now supports both course-level and lecture-level tracking.';
        COMMENT ON COLUMN course_resource_analytics.lecture_id IS 'Optional reference to specific lecture for lecture-level analytics';
    END IF;
END $$;

-- Create a view that shows resource distribution across lectures and courses
CREATE OR REPLACE VIEW resource_distribution_view AS
SELECT
    c.id as course_id,
    c.title as course_title,
    l.id as lecture_id,
    l.title as lecture_title,
    l.order_number,
    COUNT(lr.id) as lecture_resource_count,
    (SELECT COUNT(*) FROM course_resources cr WHERE cr.course_id = c.id) as course_resource_count,
    ARRAY_AGG(DISTINCT lr.resource_type) FILTER (WHERE lr.resource_type IS NOT NULL) as lecture_resource_types
FROM courses c
LEFT JOIN lectures l ON l.course_id = c.id
LEFT JOIN lecture_resources lr ON lr.lecture_id = l.id
WHERE c.deleted_at IS NULL AND (l.deleted_at IS NULL OR l.deleted_at IS NULL)
GROUP BY c.id, c.title, l.id, l.title, l.order_number
ORDER BY c.title, l.order_number;

COMMENT ON VIEW resource_distribution_view IS 'Shows how resources are distributed across courses and lectures';

-- Create migration log entry (if table exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'migration_logs') THEN
        INSERT INTO migration_logs (migration_name, applied_at, description)
        VALUES (
            '025_migrate_resources_to_lectures',
            NOW(),
            'Created lecture_resources table and migrated resource attachment from course-level to lecture-level'
        ) ON CONFLICT DO NOTHING;
    END IF;
END $$;
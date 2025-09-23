-- Create course_products table for Lemon Squeezy course-to-variant mappings
CREATE TABLE IF NOT EXISTS course_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL,
    variant_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_course_product UNIQUE (course_id)
);

-- Create indexes for performance
CREATE INDEX idx_course_products_course_id ON course_products(course_id);
CREATE INDEX idx_course_products_variant_id ON course_products(variant_id);
CREATE INDEX idx_course_products_price ON course_products(price);

-- Add foreign key constraint if courses table exists
-- This will be ignored if the constraint already exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'courses') THEN
        ALTER TABLE course_products
        ADD CONSTRAINT fk_course_products_course
        FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;
    END IF;
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Add comment for documentation
COMMENT ON TABLE course_products IS 'Maps courses to Lemon Squeezy product variants for payment processing';
COMMENT ON COLUMN course_products.course_id IS 'Reference to the course in the courses table';
COMMENT ON COLUMN course_products.variant_id IS 'Lemon Squeezy variant ID for this course';
COMMENT ON COLUMN course_products.title IS 'Product title as displayed in Lemon Squeezy';
COMMENT ON COLUMN course_products.description IS 'Product description for Lemon Squeezy checkout';
COMMENT ON COLUMN course_products.price IS 'Course price in the specified currency';
COMMENT ON COLUMN course_products.currency IS 'Currency code (USD, EUR, etc.)';
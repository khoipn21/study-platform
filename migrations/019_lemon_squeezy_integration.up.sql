-- Add Lemon Squeezy specific columns to payment-related tables

-- Update courses table to better support paid/free model
ALTER TABLE courses ADD COLUMN IF NOT EXISTS lemon_squeezy_product_id VARCHAR(100);
ALTER TABLE courses ADD COLUMN IF NOT EXISTS lemon_squeezy_variant_id VARCHAR(100);
ALTER TABLE courses ADD COLUMN IF NOT EXISTS is_paid BOOLEAN DEFAULT FALSE;

-- Update payment_methods table for Lemon Squeezy integration
ALTER TABLE payment_methods DROP CONSTRAINT IF EXISTS payment_methods_provider_check;
ALTER TABLE payment_methods ADD CONSTRAINT payment_methods_provider_check
    CHECK (provider IN ('lemonsqueezy'));

-- Add Lemon Squeezy specific columns to transactions
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS lemon_squeezy_order_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS lemon_squeezy_checkout_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS webhook_event_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS custom_data JSONB;

-- Create table for Lemon Squeezy webhook events
CREATE TABLE IF NOT EXISTS lemon_squeezy_webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id VARCHAR(100) UNIQUE NOT NULL,
    event_name VARCHAR(50) NOT NULL,
    processed_at TIMESTAMP,
    payload JSONB NOT NULL,
    signature VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create table for Lemon Squeezy products (cache)
CREATE TABLE IF NOT EXISTS lemon_squeezy_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lemon_squeezy_product_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    store_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create table for Lemon Squeezy variants (cache)
CREATE TABLE IF NOT EXISTS lemon_squeezy_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lemon_squeezy_variant_id VARCHAR(100) UNIQUE NOT NULL,
    lemon_squeezy_product_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_courses_lemon_squeezy_product_id ON courses(lemon_squeezy_product_id);
CREATE INDEX IF NOT EXISTS idx_courses_lemon_squeezy_variant_id ON courses(lemon_squeezy_variant_id);
CREATE INDEX IF NOT EXISTS idx_courses_is_paid ON courses(is_paid);

CREATE INDEX IF NOT EXISTS idx_transactions_lemon_squeezy_order_id ON transactions(lemon_squeezy_order_id);
CREATE INDEX IF NOT EXISTS idx_transactions_lemon_squeezy_checkout_id ON transactions(lemon_squeezy_checkout_id);
CREATE INDEX IF NOT EXISTS idx_transactions_webhook_event_id ON transactions(webhook_event_id);

CREATE INDEX IF NOT EXISTS idx_webhook_events_event_id ON lemon_squeezy_webhook_events(event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_event_name ON lemon_squeezy_webhook_events(event_name);
CREATE INDEX IF NOT EXISTS idx_webhook_events_processed_at ON lemon_squeezy_webhook_events(processed_at);

CREATE INDEX IF NOT EXISTS idx_lemon_squeezy_products_product_id ON lemon_squeezy_products(lemon_squeezy_product_id);
CREATE INDEX IF NOT EXISTS idx_lemon_squeezy_variants_variant_id ON lemon_squeezy_variants(lemon_squeezy_variant_id);
CREATE INDEX IF NOT EXISTS idx_lemon_squeezy_variants_product_id ON lemon_squeezy_variants(lemon_squeezy_product_id);

-- Add foreign key constraint for variants
ALTER TABLE lemon_squeezy_variants
ADD CONSTRAINT fk_lemon_squeezy_variants_product
FOREIGN KEY (lemon_squeezy_product_id)
REFERENCES lemon_squeezy_products(lemon_squeezy_product_id)
ON DELETE CASCADE;

-- Update existing courses to be free by default (unless they have a price > 0)
UPDATE courses
SET is_paid = CASE
    WHEN price > 0 THEN TRUE
    ELSE FALSE
END,
is_free = CASE
    WHEN price > 0 THEN FALSE
    ELSE TRUE
END
WHERE is_paid IS NULL;

-- Remove old payment provider constraints and data (Stripe, PayPal)
-- Clean up old payment methods that are not Lemon Squeezy
DELETE FROM payment_methods WHERE provider NOT IN ('lemonsqueezy');

-- Clean up old provider constants from transactions if needed
UPDATE transactions
SET status = 'cancelled'
WHERE payment_method_id IN (
    SELECT id FROM payment_methods WHERE provider NOT IN ('lemonsqueezy')
);
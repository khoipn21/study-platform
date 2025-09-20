-- Rollback Lemon Squeezy integration

-- Drop Lemon Squeezy specific indexes
DROP INDEX IF EXISTS idx_courses_lemon_squeezy_product_id;
DROP INDEX IF EXISTS idx_courses_lemon_squeezy_variant_id;
DROP INDEX IF EXISTS idx_transactions_lemon_squeezy_order_id;
DROP INDEX IF EXISTS idx_transactions_lemon_squeezy_checkout_id;
DROP INDEX IF EXISTS idx_transactions_webhook_event_id;
DROP INDEX IF EXISTS idx_webhook_events_event_id;
DROP INDEX IF EXISTS idx_webhook_events_event_name;
DROP INDEX IF EXISTS idx_webhook_events_processed_at;
DROP INDEX IF EXISTS idx_lemon_squeezy_products_product_id;
DROP INDEX IF EXISTS idx_lemon_squeezy_variants_variant_id;
DROP INDEX IF EXISTS idx_lemon_squeezy_variants_product_id;

-- Drop Lemon Squeezy specific tables
DROP TABLE IF EXISTS lemon_squeezy_variants;
DROP TABLE IF EXISTS lemon_squeezy_products;
DROP TABLE IF EXISTS lemon_squeezy_webhook_events;

-- Remove Lemon Squeezy specific columns from transactions
ALTER TABLE transactions DROP COLUMN IF EXISTS lemon_squeezy_order_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS lemon_squeezy_checkout_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS webhook_event_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS custom_data;

-- Remove Lemon Squeezy specific columns from courses
ALTER TABLE courses DROP COLUMN IF EXISTS lemon_squeezy_product_id;
ALTER TABLE courses DROP COLUMN IF EXISTS lemon_squeezy_variant_id;
ALTER TABLE courses DROP COLUMN IF EXISTS is_paid;

-- Restore original payment method provider constraint
ALTER TABLE payment_methods DROP CONSTRAINT IF EXISTS payment_methods_provider_check;
ALTER TABLE payment_methods ADD CONSTRAINT payment_methods_provider_check
    CHECK (provider IN ('stripe', 'paypal'));
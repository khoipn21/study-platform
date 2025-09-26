-- Remove Stripe-specific columns from transactions table

-- Drop constraints first
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS uk_transactions_stripe_payment_intent_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_stripe_payment_intent_id;
DROP INDEX IF EXISTS idx_transactions_stripe_customer_id;
DROP INDEX IF EXISTS idx_transactions_stripe_charge_id;
DROP INDEX IF EXISTS idx_transactions_lemon_squeezy_order_id;
DROP INDEX IF EXISTS idx_transactions_lemon_squeezy_checkout_id;

-- Drop columns
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_payment_intent_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_customer_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_charge_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_session_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_invoice_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS stripe_subscription_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS lemon_squeezy_order_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS lemon_squeezy_checkout_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS webhook_event_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS custom_data;
-- Add Stripe-specific columns to transactions table
-- This migration adds missing Stripe columns that are required by the repository code

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_payment_intent_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_customer_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_charge_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_session_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_invoice_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS stripe_subscription_id VARCHAR(100);

-- Add LemonSqueezy columns that are referenced but missing
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS lemon_squeezy_order_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS lemon_squeezy_checkout_id VARCHAR(100);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS webhook_event_id UUID;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS custom_data JSONB DEFAULT '{}';

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transactions_stripe_payment_intent_id ON transactions(stripe_payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_transactions_stripe_customer_id ON transactions(stripe_customer_id);
CREATE INDEX IF NOT EXISTS idx_transactions_stripe_charge_id ON transactions(stripe_charge_id);
CREATE INDEX IF NOT EXISTS idx_transactions_lemon_squeezy_order_id ON transactions(lemon_squeezy_order_id);
CREATE INDEX IF NOT EXISTS idx_transactions_lemon_squeezy_checkout_id ON transactions(lemon_squeezy_checkout_id);

-- Add unique constraints where appropriate
ALTER TABLE transactions ADD CONSTRAINT uk_transactions_stripe_payment_intent_id
    UNIQUE (stripe_payment_intent_id) DEFERRABLE INITIALLY DEFERRED;

COMMENT ON COLUMN transactions.stripe_payment_intent_id IS 'Stripe payment intent ID for tracking payments';
COMMENT ON COLUMN transactions.stripe_customer_id IS 'Stripe customer ID associated with the transaction';
COMMENT ON COLUMN transactions.stripe_charge_id IS 'Stripe charge ID for completed payments';
COMMENT ON COLUMN transactions.lemon_squeezy_order_id IS 'LemonSqueezy order ID for alternative payment processor';
COMMENT ON COLUMN transactions.lemon_squeezy_checkout_id IS 'LemonSqueezy checkout ID for alternative payment processor';
COMMENT ON COLUMN transactions.custom_data IS 'Additional payment metadata in JSON format';
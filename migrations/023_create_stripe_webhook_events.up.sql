-- Create Stripe webhook events table
-- This table stores Stripe webhook events for processing and deduplication

CREATE TABLE IF NOT EXISTS stripe_webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_event_id VARCHAR(100) NOT NULL UNIQUE,
    event_type VARCHAR(100) NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    processing_attempts INTEGER DEFAULT 0,
    event_data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    error_message TEXT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_stripe_event_id ON stripe_webhook_events(stripe_event_id);
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_event_type ON stripe_webhook_events(event_type);
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_processed ON stripe_webhook_events(processed);
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_created_at ON stripe_webhook_events(created_at);
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_processing_attempts ON stripe_webhook_events(processing_attempts);

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_stripe_webhook_events_unprocessed ON stripe_webhook_events(processed, processing_attempts, created_at)
    WHERE processed = FALSE;

COMMENT ON TABLE stripe_webhook_events IS 'Stripe webhook events storage for processing and deduplication';
COMMENT ON COLUMN stripe_webhook_events.stripe_event_id IS 'Unique event ID from Stripe';
COMMENT ON COLUMN stripe_webhook_events.event_type IS 'Type of Stripe event (e.g., payment_intent.succeeded)';
COMMENT ON COLUMN stripe_webhook_events.processed IS 'Whether the event has been successfully processed';
COMMENT ON COLUMN stripe_webhook_events.processing_attempts IS 'Number of times processing was attempted';
COMMENT ON COLUMN stripe_webhook_events.event_data IS 'Full JSON payload from Stripe webhook';
COMMENT ON COLUMN stripe_webhook_events.error_message IS 'Error message from last failed processing attempt';
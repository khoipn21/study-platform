-- Drop indexes
DROP INDEX IF EXISTS idx_network_metrics_timestamp;
DROP INDEX IF EXISTS idx_network_metrics_session;

-- Drop network_metrics table
DROP TABLE IF EXISTS network_metrics;
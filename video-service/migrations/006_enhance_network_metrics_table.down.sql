-- Rollback enhanced network metrics table
DROP TRIGGER IF EXISTS trigger_update_network_analytics ON network_metrics;
DROP FUNCTION IF EXISTS update_network_analytics_daily();
DROP FUNCTION IF EXISTS refresh_network_dashboard();
DROP MATERIALIZED VIEW IF EXISTS network_monitoring_dashboard;
DROP TABLE IF EXISTS bandwidth_tests;
DROP TABLE IF EXISTS network_analytics_daily;
DROP TABLE IF EXISTS adaptive_streaming_rules;
DROP TABLE IF EXISTS network_events;

-- Remove added columns from network_metrics table
ALTER TABLE network_metrics
DROP COLUMN IF EXISTS device_type,
DROP COLUMN IF EXISTS screen_resolution,
DROP COLUMN IF EXISTS user_agent,
DROP COLUMN IF EXISTS geographic_location,
DROP COLUMN IF EXISTS isp_info,
DROP COLUMN IF EXISTS network_stability_score,
DROP COLUMN IF EXISTS jitter_ms,
DROP COLUMN IF EXISTS download_speed_mbps,
DROP COLUMN IF EXISTS upload_speed_mbps,
DROP COLUMN IF EXISTS cpu_usage_percent,
DROP COLUMN IF EXISTS memory_usage_percent,
DROP COLUMN IF EXISTS battery_level,
DROP COLUMN IF EXISTS thermal_state;

-- Remove enhanced indexes
DROP INDEX IF EXISTS idx_network_metrics_user_timestamp;
DROP INDEX IF EXISTS idx_network_metrics_quality_score;
DROP INDEX IF EXISTS idx_network_metrics_connection_type;
DROP INDEX IF EXISTS idx_network_metrics_device_type;
-- Enhance network metrics table for comprehensive monitoring
ALTER TABLE network_metrics
ADD COLUMN IF NOT EXISTS device_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS screen_resolution VARCHAR(20),
ADD COLUMN IF NOT EXISTS user_agent TEXT,
ADD COLUMN IF NOT EXISTS geographic_location JSONB,
ADD COLUMN IF NOT EXISTS isp_info JSONB,
ADD COLUMN IF NOT EXISTS network_stability_score INTEGER DEFAULT 5,
ADD COLUMN IF NOT EXISTS jitter_ms INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS download_speed_mbps DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS upload_speed_mbps DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS cpu_usage_percent DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS memory_usage_percent DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS battery_level INTEGER,
ADD COLUMN IF NOT EXISTS thermal_state VARCHAR(20);

-- Create enhanced indexes for network metrics
CREATE INDEX IF NOT EXISTS idx_network_metrics_user_timestamp ON network_metrics(user_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_network_metrics_quality_score ON network_metrics(quality_score);
CREATE INDEX IF NOT EXISTS idx_network_metrics_connection_type ON network_metrics(connection_type);
CREATE INDEX IF NOT EXISTS idx_network_metrics_device_type ON network_metrics(device_type);

-- Create table for real-time network events
CREATE TABLE IF NOT EXISTS network_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'quality_change', 'buffer_event', 'connection_lost', 'recovery'
    event_data JSONB,
    severity VARCHAR(20) DEFAULT 'info', -- 'info', 'warning', 'error', 'critical'
    timestamp TIMESTAMP DEFAULT NOW(),
    resolved BOOLEAN DEFAULT FALSE,
    resolution_timestamp TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for network events
CREATE INDEX idx_network_events_session ON network_events(session_id);
CREATE INDEX idx_network_events_user ON network_events(user_id);
CREATE INDEX idx_network_events_type ON network_events(event_type);
CREATE INDEX idx_network_events_timestamp ON network_events(timestamp);
CREATE INDEX idx_network_events_severity ON network_events(severity);
CREATE INDEX idx_network_events_unresolved ON network_events(resolved) WHERE resolved = FALSE;

-- Create table for adaptive streaming rules
CREATE TABLE IF NOT EXISTS adaptive_streaming_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_name VARCHAR(100) NOT NULL,
    condition_bandwidth_min DECIMAL(10,2),
    condition_bandwidth_max DECIMAL(10,2),
    condition_latency_max INTEGER,
    condition_packet_loss_max DECIMAL(5,2),
    condition_connection_types TEXT[], -- Array of connection types
    condition_device_types TEXT[], -- Array of device types
    recommended_quality VARCHAR(10) NOT NULL,
    buffer_target_seconds INTEGER DEFAULT 10,
    preload_enabled BOOLEAN DEFAULT TRUE,
    active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 50, -- Higher numbers = higher priority
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for adaptive streaming rules
CREATE INDEX idx_adaptive_rules_active ON adaptive_streaming_rules(active);
CREATE INDEX idx_adaptive_rules_priority ON adaptive_streaming_rules(priority DESC);

-- Insert default adaptive streaming rules
INSERT INTO adaptive_streaming_rules (rule_name, condition_bandwidth_min, condition_bandwidth_max, condition_latency_max, condition_packet_loss_max, condition_connection_types, recommended_quality, buffer_target_seconds, priority)
VALUES
    ('Ultra High Quality', 15.0, NULL, 50, 0.005, ARRAY['5g', 'ethernet'], '1080p', 8, 90),
    ('High Quality WiFi', 8.0, 15.0, 100, 0.01, ARRAY['wifi'], '720p', 10, 80),
    ('Standard 4G', 3.0, 8.0, 150, 0.02, ARRAY['4g'], '480p', 12, 70),
    ('Low Bandwidth', 1.0, 3.0, 200, 0.03, NULL, '360p', 15, 60),
    ('Emergency Mode', 0.5, 1.0, 500, 0.05, NULL, '240p', 20, 50),
    ('Excellent Ethernet', 25.0, NULL, 30, 0.001, ARRAY['ethernet'], '1080p', 6, 95)
ON CONFLICT DO NOTHING;

-- Create table for network analytics aggregates
CREATE TABLE IF NOT EXISTS network_analytics_daily (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    user_id UUID,
    video_id UUID,
    session_count INTEGER DEFAULT 0,
    avg_bandwidth_mbps DECIMAL(10,2),
    avg_latency_ms DECIMAL(8,2),
    avg_packet_loss DECIMAL(5,4),
    quality_distribution JSONB, -- {"720p": 60, "1080p": 30, "480p": 10}
    connection_type_distribution JSONB,
    device_type_distribution JSONB,
    total_quality_changes INTEGER DEFAULT 0,
    total_buffer_events INTEGER DEFAULT 0,
    total_network_interruptions INTEGER DEFAULT 0,
    avg_stability_score DECIMAL(4,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(date, user_id, video_id)
);

-- Indexes for network analytics
CREATE INDEX idx_network_analytics_date ON network_analytics_daily(date);
CREATE INDEX idx_network_analytics_user ON network_analytics_daily(user_id);
CREATE INDEX idx_network_analytics_video ON network_analytics_daily(video_id);

-- Create table for bandwidth test results
CREATE TABLE IF NOT EXISTS bandwidth_tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    test_type VARCHAR(20) NOT NULL, -- 'initial', 'periodic', 'triggered'
    download_mbps DECIMAL(10,2),
    upload_mbps DECIMAL(10,2),
    latency_ms INTEGER,
    jitter_ms INTEGER,
    packet_loss_percent DECIMAL(5,4),
    test_duration_seconds INTEGER,
    test_server_location VARCHAR(100),
    confidence_score DECIMAL(3,2), -- 0.0 to 1.0
    timestamp TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for bandwidth tests
CREATE INDEX idx_bandwidth_tests_session ON bandwidth_tests(session_id);
CREATE INDEX idx_bandwidth_tests_user ON bandwidth_tests(user_id);
CREATE INDEX idx_bandwidth_tests_timestamp ON bandwidth_tests(timestamp);
CREATE INDEX idx_bandwidth_tests_type ON bandwidth_tests(test_type);

-- Create materialized view for real-time network monitoring dashboard
CREATE MATERIALIZED VIEW IF NOT EXISTS network_monitoring_dashboard AS
SELECT
    nm.session_id,
    nm.user_id,
    vs.video_id,
    nm.timestamp,
    nm.bandwidth_mbps,
    nm.latency_ms,
    nm.packet_loss_percent,
    nm.connection_type,
    nm.device_type,
    nm.quality_score,
    nm.recommended_quality,
    nm.buffer_health_seconds,
    nm.network_stability_score,
    vs.current_quality,
    vs.current_time_seconds,
    CASE
        WHEN nm.timestamp > NOW() - INTERVAL '5 minutes' THEN 'active'
        WHEN nm.timestamp > NOW() - INTERVAL '30 minutes' THEN 'recent'
        ELSE 'inactive'
    END as session_status,
    EXTRACT(EPOCH FROM (NOW() - nm.timestamp)) as seconds_since_last_update
FROM network_metrics nm
LEFT JOIN viewing_sessions vs ON nm.session_id = vs.session_id
WHERE nm.timestamp > NOW() - INTERVAL '4 hours';

-- Index for materialized view
CREATE INDEX IF NOT EXISTS idx_network_dashboard_session_status ON network_monitoring_dashboard(session_status);
CREATE INDEX IF NOT EXISTS idx_network_dashboard_video_id ON network_monitoring_dashboard(video_id);

-- Function to refresh the materialized view
CREATE OR REPLACE FUNCTION refresh_network_dashboard()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY network_monitoring_dashboard;
END;
$$ LANGUAGE plpgsql;

-- Create trigger function for automatic network analytics aggregation
CREATE OR REPLACE FUNCTION update_network_analytics_daily()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO network_analytics_daily (
        date, user_id, video_id, session_count, avg_bandwidth_mbps,
        avg_latency_ms, avg_packet_loss, total_quality_changes
    )
    VALUES (
        CURRENT_DATE, NEW.user_id,
        (SELECT video_id FROM viewing_sessions WHERE session_id = NEW.session_id LIMIT 1),
        1, NEW.bandwidth_mbps, NEW.latency_ms, NEW.packet_loss_percent, 0
    )
    ON CONFLICT (date, user_id, video_id)
    DO UPDATE SET
        session_count = network_analytics_daily.session_count + 1,
        avg_bandwidth_mbps = (network_analytics_daily.avg_bandwidth_mbps * network_analytics_daily.session_count + NEW.bandwidth_mbps) / (network_analytics_daily.session_count + 1),
        avg_latency_ms = (network_analytics_daily.avg_latency_ms * network_analytics_daily.session_count + NEW.latency_ms) / (network_analytics_daily.session_count + 1),
        avg_packet_loss = (network_analytics_daily.avg_packet_loss * network_analytics_daily.session_count + NEW.packet_loss_percent) / (network_analytics_daily.session_count + 1),
        updated_at = NOW();

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for network metrics
DROP TRIGGER IF EXISTS trigger_update_network_analytics ON network_metrics;
CREATE TRIGGER trigger_update_network_analytics
    AFTER INSERT ON network_metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_network_analytics_daily();
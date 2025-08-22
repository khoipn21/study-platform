-- Network quality metrics
CREATE TABLE network_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW(),
    bandwidth_mbps DECIMAL(10,2),
    latency_ms INTEGER,
    packet_loss_percent DECIMAL(5,2),
    connection_type VARCHAR(20), -- 'wifi', '4g', '5g', 'ethernet'
    quality_score INTEGER, -- 1-10 scale
    recommended_quality VARCHAR(10),
    buffer_health_seconds INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for network_metrics table
CREATE INDEX idx_network_metrics_session ON network_metrics(session_id);
CREATE INDEX idx_network_metrics_timestamp ON network_metrics(timestamp);
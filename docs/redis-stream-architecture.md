# Redis Stream Architecture for Video Streaming Platform

## Overview

This document outlines the Redis Stream architecture designed for real-time network quality monitoring, adaptive video streaming, and user engagement tracking. The system uses Redis Streams to handle high-throughput, time-ordered event processing with guaranteed delivery and consumer group management.

## Redis Stream Architecture

### Stream Design Patterns

#### 1. Network Quality Monitoring Stream
```
Stream Key: network:quality:{user_id}
Purpose: Track network conditions in real-time
TTL: 24 hours (configurable)
```

**Event Structure:**
```json
{
  "event_id": "1694982000000-0",
  "session_id": "sess_abc123",
  "user_id": "user_xyz789",
  "video_id": "video_def456",
  "timestamp": "2024-09-17T10:30:00Z",
  "bandwidth_mbps": 15.2,
  "latency_ms": 45,
  "packet_loss_percent": 0.1,
  "connection_type": "wifi",
  "quality_score": 8,
  "buffer_health_seconds": 12,
  "current_quality": "1080p",
  "recommended_quality": "1080p",
  "jitter_ms": 5,
  "throughput_mbps": 14.8
}
```

#### 2. Video Quality Adaptation Stream
```
Stream Key: video:quality:changes
Purpose: Track and coordinate quality changes across all sessions
TTL: 7 days
```

**Event Structure:**
```json
{
  "event_id": "1694982000000-1",
  "session_id": "sess_abc123",
  "video_id": "video_def456",
  "user_id": "user_xyz789",
  "from_quality": "1080p",
  "to_quality": "720p",
  "reason": "bandwidth_drop",
  "trigger_metric": "bandwidth",
  "metric_value": 8.5,
  "auto_switch": true,
  "timestamp": "2024-09-17T10:30:15Z",
  "confidence": 0.92
}
```

#### 3. User Engagement Stream
```
Stream Key: engagement:events:{video_id}
Purpose: Track user interactions and video engagement
TTL: 30 days
```

**Event Structure:**
```json
{
  "event_id": "1694982000000-2",
  "session_id": "sess_abc123",
  "user_id": "user_xyz789",
  "video_id": "video_def456",
  "event_type": "seek",
  "timestamp": "2024-09-17T10:30:20Z",
  "video_timestamp_seconds": 125,
  "event_data": {
    "from_position": 120,
    "to_position": 180,
    "seek_type": "manual"
  }
}
```

#### 4. Real-time Analytics Stream
```
Stream Key: analytics:realtime
Purpose: Aggregate real-time metrics for dashboards
TTL: 1 hour
```

**Event Structure:**
```json
{
  "event_id": "1694982000000-3",
  "timestamp": "2024-09-17T10:30:00Z",
  "metric_type": "concurrent_viewers",
  "video_id": "video_def456",
  "value": 1247,
  "additional_data": {
    "quality_distribution": {
      "360p": 123,
      "720p": 456,
      "1080p": 668
    },
    "geographic_distribution": {
      "US": 747,
      "EU": 312,
      "AS": 188
    }
  }
}
```

### Consumer Groups Configuration

#### 1. Network Intelligence Service
```
Consumer Group: network-intelligence
Stream: network:quality:*
Purpose: Process network metrics and generate recommendations
Consumers: 3-5 instances (auto-scaling)
Processing: Real-time (<100ms)
```

#### 2. Quality Adaptation Service
```
Consumer Group: quality-adapter
Stream: video:quality:changes
Purpose: Coordinate quality changes and update streaming parameters
Consumers: 2-3 instances
Processing: Real-time (<50ms)
```

#### 3. Analytics Aggregator
```
Consumer Group: analytics-processor
Stream: engagement:events:*, analytics:realtime
Purpose: Aggregate metrics for reporting and dashboards
Consumers: 2-4 instances
Processing: Near real-time (<1s)
```

#### 4. Buffering Optimizer
```
Consumer Group: buffer-optimizer
Stream: network:quality:*
Purpose: Optimize preloading and buffer management
Consumers: 2-3 instances
Processing: Real-time (<100ms)
```

## Stream Processing Patterns

### 1. Network Quality Assessment Algorithm

```go
// Pseudo-code for network quality processing
func ProcessNetworkMetrics(event NetworkEvent) {
    // Calculate quality score based on multiple factors
    qualityScore := calculateQualityScore(
        event.BandwidthMbps,
        event.LatencyMs,
        event.PacketLossPercent,
        event.JitterMs,
    )

    // Determine recommended quality
    recommendedQuality := determineOptimalQuality(
        qualityScore,
        event.CurrentQuality,
        event.BufferHealthSeconds,
    )

    // Check if quality change is needed
    if shouldChangeQuality(event.CurrentQuality, recommendedQuality) {
        publishQualityChangeEvent(event.SessionID, recommendedQuality)
    }

    // Update buffer strategy
    updateBufferStrategy(event.SessionID, qualityScore)
}
```

### 2. Adaptive Streaming Decision Engine

```go
func calculateQualityScore(bandwidth, latency, packetLoss, jitter float64) int {
    score := 10.0

    // Bandwidth scoring (40% weight)
    if bandwidth < 2.0 {
        score -= 4.0
    } else if bandwidth < 5.0 {
        score -= 2.0
    } else if bandwidth < 10.0 {
        score -= 1.0
    }

    // Latency scoring (30% weight)
    if latency > 200 {
        score -= 3.0
    } else if latency > 100 {
        score -= 1.5
    } else if latency > 50 {
        score -= 0.5
    }

    // Packet loss scoring (20% weight)
    if packetLoss > 5.0 {
        score -= 2.0
    } else if packetLoss > 1.0 {
        score -= 1.0
    } else if packetLoss > 0.5 {
        score -= 0.5
    }

    // Jitter scoring (10% weight)
    if jitter > 50 {
        score -= 1.0
    } else if jitter > 20 {
        score -= 0.5
    }

    return int(math.Max(1, math.Min(10, score)))
}
```

### 3. Quality Mapping Strategy

```go
var qualityMap = map[int][]string{
    1:  {"360p"},                    // Very poor connection
    2:  {"360p"},                    // Poor connection
    3:  {"360p", "480p"},           // Below average
    4:  {"480p"},                    // Below average
    5:  {"480p", "720p"},           // Average
    6:  {"720p"},                    // Good
    7:  {"720p", "1080p"},          // Very good
    8:  {"1080p"},                   // Excellent
    9:  {"1080p", "1440p"},         // Excellent+
    10: {"1080p", "1440p", "2160p"}, // Perfect
}
```

## Redis Configuration and Optimization

### 1. Memory Management
```redis
# Configure memory policy for streams
maxmemory-policy allkeys-lru
maxmemory 8gb

# Stream-specific settings
stream-node-max-bytes 4096
stream-node-max-entries 100
```

### 2. Persistence Configuration
```redis
# AOF for better durability of streams
appendonly yes
appendfsync everysec
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
```

### 3. Stream Trimming Strategy
```redis
# Automatic trimming to manage memory usage
XTRIM network:quality:* MAXLEN ~ 10000
XTRIM video:quality:changes MAXLEN ~ 50000
XTRIM engagement:events:* MAXLEN ~ 100000
XTRIM analytics:realtime MAXLEN ~ 1000
```

## WebSocket Integration

### 1. Real-time Communication Flow

```
Client (Video Player)
    ↓ (WebSocket)
Video Service
    ↓ (Redis Stream)
Network Intelligence Service
    ↓ (Redis Stream)
Quality Adaptation Service
    ↓ (WebSocket)
Client (Video Player)
```

### 2. WebSocket Message Types

#### Network Status Update (Client → Server)
```json
{
  "type": "network_status",
  "data": {
    "session_id": "sess_abc123",
    "bandwidth_mbps": 15.2,
    "latency_ms": 45,
    "buffer_health": 12,
    "current_quality": "1080p",
    "timestamp": "2024-09-17T10:30:00Z"
  }
}
```

#### Quality Recommendation (Server → Client)
```json
{
  "type": "quality_recommendation",
  "data": {
    "recommended_quality": "720p",
    "reason": "bandwidth_optimization",
    "confidence": 0.92,
    "should_preload": true,
    "buffer_target": 15
  }
}
```

#### Preload Instruction (Server → Client)
```json
{
  "type": "preload_instruction",
  "data": {
    "segments": [
      "segment_001.m4s",
      "segment_002.m4s"
    ],
    "priority": "high",
    "quality": "720p"
  }
}
```

## Performance Monitoring and Metrics

### 1. Stream Health Metrics
- **Throughput**: Events per second per stream
- **Latency**: Time from event creation to processing
- **Consumer Lag**: Delay in consumer group processing
- **Memory Usage**: Stream memory consumption
- **Error Rate**: Failed event processing percentage

### 2. Key Performance Indicators (KPIs)
- **Average Quality Score**: Overall network quality across users
- **Quality Change Frequency**: Adaptations per viewing session
- **Buffering Incident Rate**: Percentage of sessions with buffering
- **Startup Time**: Time to first video frame
- **Completion Rate**: Percentage of videos watched to completion

### 3. Monitoring Queries
```redis
# Get stream information
XINFO STREAM network:quality:user123

# Check consumer group lag
XINFO GROUPS video:quality:changes

# Monitor stream length
XLEN engagement:events:video456

# Get latest entries
XREVRANGE analytics:realtime + - COUNT 10
```

## Scalability Considerations

### 1. Horizontal Scaling
- **Stream Partitioning**: Partition streams by user_id hash
- **Consumer Scaling**: Auto-scale consumers based on lag
- **Redis Clustering**: Use Redis Cluster for distributed streams
- **Load Balancing**: Balance WebSocket connections across instances

### 2. Data Lifecycle Management
- **Stream TTL**: Automatic expiration of old stream data
- **Archival**: Move historical data to PostgreSQL
- **Compression**: Compress older stream entries
- **Cleanup**: Regular cleanup of inactive streams

### 3. Fault Tolerance
- **Consumer Recovery**: Automatic restart of failed consumers
- **Data Replication**: Redis replication for high availability
- **Circuit Breakers**: Prevent cascade failures
- **Graceful Degradation**: Fallback to basic quality when streams fail

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Set up Redis Streams configuration
- [ ] Implement basic stream producers
- [ ] Create consumer group management
- [ ] Set up WebSocket handling

### Phase 2: Network Intelligence
- [ ] Implement network quality assessment
- [ ] Create quality recommendation engine
- [ ] Set up adaptive streaming logic
- [ ] Add buffering optimization

### Phase 3: Analytics and Monitoring
- [ ] Create real-time analytics aggregation
- [ ] Set up monitoring dashboards
- [ ] Implement alert systems
- [ ] Add performance metrics

### Phase 4: Optimization and Scaling
- [ ] Optimize stream performance
- [ ] Implement auto-scaling
- [ ] Add fault tolerance mechanisms
- [ ] Create operational runbooks

This Redis Stream architecture provides a robust foundation for real-time video streaming optimization with network intelligence, ensuring optimal user experience across varying network conditions.
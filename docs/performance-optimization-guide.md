# Performance Optimization Guide for Video Streaming Platform

## Overview

This guide provides comprehensive performance optimization strategies for the video streaming platform, covering database optimization, Redis stream performance, video delivery optimization, and scalability patterns for smooth video streaming experiences.

## Database Performance Optimization

### 1. PostgreSQL Configuration for Video Streaming

#### Connection Pooling Configuration
```postgresql
# postgresql.conf optimizations for video platform
max_connections = 500
shared_buffers = 2GB
effective_cache_size = 6GB
work_mem = 16MB
maintenance_work_mem = 256MB

# Write-ahead logging optimization
wal_buffers = 64MB
checkpoint_segments = 64
checkpoint_completion_target = 0.9

# Background writer optimization
bgwriter_delay = 50ms
bgwriter_lru_maxpages = 1000
bgwriter_lru_multiplier = 10.0
```

#### Index Optimization Strategies

**Time-Series Data Optimization**:
```sql
-- Partition network_metrics table by timestamp for better performance
CREATE TABLE network_metrics_y2024m09 PARTITION OF network_metrics
    FOR VALUES FROM ('2024-09-01') TO ('2024-10-01');

-- Composite indexes for common query patterns
CREATE INDEX CONCURRENTLY idx_network_metrics_session_timestamp
    ON network_metrics (session_id, timestamp DESC);

CREATE INDEX CONCURRENTLY idx_viewing_sessions_user_video_started
    ON viewing_sessions (user_id, video_id, started_at DESC);

-- Partial indexes for active sessions
CREATE INDEX CONCURRENTLY idx_viewing_sessions_active
    ON viewing_sessions (last_heartbeat)
    WHERE completed = false;
```

**Video Analytics Optimization**:
```sql
-- Covering index for video analytics queries
CREATE INDEX CONCURRENTLY idx_video_analytics_covering
    ON video_analytics (video_id, date DESC)
    INCLUDE (total_views, unique_viewers, avg_watch_time_seconds);

-- Expression index for quality distribution queries
CREATE INDEX CONCURRENTLY idx_video_analytics_quality_jsonb
    ON video_analytics USING GIN (quality_distribution);
```

### 2. Read Replica Strategy

```yaml
# Database read replica configuration
database:
  write_primary:
    host: "primary-db.internal"
    max_connections: 200
    pool_size: 50

  read_replicas:
    analytics:
      host: "analytics-replica.internal"
      max_connections: 100
      pool_size: 25
      purpose: "video_analytics,reporting"

    sessions:
      host: "sessions-replica.internal"
      max_connections: 150
      pool_size: 35
      purpose: "viewing_sessions,network_metrics"
```

### 3. Database Query Optimization

#### Optimized Analytics Queries
```sql
-- Materialized view for video performance dashboard
CREATE MATERIALIZED VIEW video_performance_hourly AS
SELECT
    date_trunc('hour', vs.started_at) as hour,
    v.id as video_id,
    v.title,
    COUNT(DISTINCT vs.user_id) as unique_viewers,
    AVG(vs.total_watch_time_seconds) as avg_watch_time,
    AVG(nm.quality_score) as avg_quality_score,
    COUNT(qc.id) as quality_changes,
    COUNT(be.id) as buffering_events
FROM videos v
LEFT JOIN viewing_sessions vs ON v.id = vs.video_id
LEFT JOIN network_metrics nm ON vs.session_id = nm.session_id
LEFT JOIN quality_changes qc ON vs.session_id = qc.session_id
LEFT JOIN buffering_events be ON vs.session_id = be.session_id
WHERE vs.started_at >= NOW() - INTERVAL '7 days'
GROUP BY 1, 2, 3;

-- Refresh strategy
SELECT cron.schedule('refresh-video-performance', '*/15 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY video_performance_hourly;');
```

#### Efficient Session Queries
```sql
-- Get active sessions with network metrics
WITH active_sessions AS (
    SELECT session_id, user_id, video_id, current_quality
    FROM viewing_sessions
    WHERE last_heartbeat > NOW() - INTERVAL '30 seconds'
      AND completed = false
),
latest_metrics AS (
    SELECT DISTINCT ON (session_id)
           session_id, quality_score, bandwidth_mbps, recommended_quality
    FROM network_metrics
    WHERE timestamp > NOW() - INTERVAL '5 minutes'
    ORDER BY session_id, timestamp DESC
)
SELECT s.*, m.quality_score, m.bandwidth_mbps, m.recommended_quality
FROM active_sessions s
LEFT JOIN latest_metrics m ON s.session_id = m.session_id;
```

## Redis Performance Optimization

### 1. Redis Configuration for Streams

```redis
# redis.conf optimizations for video streaming
maxmemory 8gb
maxmemory-policy allkeys-lru

# Stream-specific configuration
stream-node-max-bytes 4096
stream-node-max-entries 100

# Persistence optimization
save 900 1
save 300 10
save 60 10000

# Network optimization
tcp-keepalive 300
timeout 0
tcp-backlog 511
```

### 2. Stream Performance Patterns

#### Efficient Stream Design
```go
// Optimized stream producer with batching
type StreamProducer struct {
    client     redis.Client
    batchSize  int
    flushTime  time.Duration
    buffer     []StreamEntry
    mutex      sync.Mutex
}

func (p *StreamProducer) PublishNetworkMetrics(userID string, metrics NetworkMetrics) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    entry := StreamEntry{
        Stream: fmt.Sprintf("network:quality:%s", userID),
        ID:     "*",
        Values: map[string]interface{}{
            "session_id":     metrics.SessionID,
            "bandwidth_mbps": metrics.BandwidthMbps,
            "latency_ms":     metrics.LatencyMs,
            "quality_score":  metrics.QualityScore,
            "timestamp":      time.Now().UnixMilli(),
        },
    }

    p.buffer = append(p.buffer, entry)

    if len(p.buffer) >= p.batchSize {
        return p.flush()
    }

    return nil
}

func (p *StreamProducer) flush() error {
    if len(p.buffer) == 0 {
        return nil
    }

    pipe := p.client.Pipeline()
    for _, entry := range p.buffer {
        pipe.XAdd(context.Background(), &redis.XAddArgs{
            Stream: entry.Stream,
            ID:     entry.ID,
            Values: entry.Values,
        })
    }

    _, err := pipe.Exec(context.Background())
    if err == nil {
        p.buffer = p.buffer[:0] // Reset buffer
    }

    return err
}
```

#### Consumer Group Optimization
```go
// High-performance stream consumer with concurrent processing
type StreamConsumer struct {
    client        redis.Client
    consumerGroup string
    consumerName  string
    workers       int
    handler       func(Entry) error
}

func (c *StreamConsumer) ConsumeWithWorkers(streams []string) error {
    workerPool := make(chan StreamEntry, c.workers*2)

    // Start worker goroutines
    for i := 0; i < c.workers; i++ {
        go c.worker(workerPool)
    }

    for {
        // Read from multiple streams
        result, err := c.client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
            Group:    c.consumerGroup,
            Consumer: c.consumerName,
            Streams:  append(streams, []string{">", ">", ">"}...),
            Count:    100,
            Block:    time.Second,
        }).Result()

        if err != nil {
            continue
        }

        // Distribute entries to workers
        for _, stream := range result {
            for _, entry := range stream.Messages {
                workerPool <- StreamEntry{
                    StreamName: stream.Stream,
                    ID:         entry.ID,
                    Values:     entry.Values,
                }
            }
        }
    }
}

func (c *StreamConsumer) worker(entries chan StreamEntry) {
    for entry := range entries {
        if err := c.handler(entry); err != nil {
            log.Printf("Error processing entry %s: %v", entry.ID, err)
            continue
        }

        // Acknowledge successful processing
        c.client.XAck(context.Background(), entry.StreamName, c.consumerGroup, entry.ID)
    }
}
```

### 3. Memory Management and Stream Trimming

```go
// Automated stream maintenance
type StreamMaintainer struct {
    client    redis.Client
    trimmer   *time.Ticker
    streams   map[string]StreamConfig
}

type StreamConfig struct {
    MaxLength int64
    TTL       time.Duration
    TrimInterval time.Duration
}

func (sm *StreamMaintainer) Start() {
    sm.trimmer = time.NewTicker(5 * time.Minute)

    go func() {
        for range sm.trimmer.C {
            sm.trimStreams()
        }
    }()
}

func (sm *StreamMaintainer) trimStreams() {
    for pattern, config := range sm.streams {
        keys, err := sm.client.Keys(context.Background(), pattern).Result()
        if err != nil {
            continue
        }

        for _, key := range keys {
            // Trim to max length
            sm.client.XTrimMaxLen(context.Background(), key, config.MaxLength)

            // Set TTL if stream is inactive
            lastEntry, err := sm.client.XRevRangeN(context.Background(), key, "+", "-", 1).Result()
            if err != nil || len(lastEntry) == 0 {
                continue
            }

            // Parse timestamp from last entry
            lastTime := sm.parseTimestamp(lastEntry[0].ID)
            if time.Since(lastTime) > config.TTL {
                sm.client.Expire(context.Background(), key, time.Hour)
            }
        }
    }
}
```

## Video Delivery Optimization

### 1. Cloudflare Stream Integration

#### Optimized Video Upload
```go
type CloudflareStreamService struct {
    client     *http.Client
    apiToken   string
    accountID  string
    baseURL    string
}

func (cs *CloudflareStreamService) OptimizedUpload(video VideoUploadRequest) (*UploadResponse, error) {
    // Create video with optimized settings
    createReq := CloudflareCreateVideoRequest{
        Meta: map[string]string{
            "name":        video.Title,
            "description": video.Description,
        },
        RequireSignedURLs: true,
        AllowedOrigins:    []string{"*.studyplatform.com"},
        Thumbnails: ThumbnailSettings{
            Unit:       "percent",
            Positions:  []int{10, 25, 50, 75, 90},
            Count:      5,
        },
        Watermark: WatermarkSettings{
            UID:     cs.getWatermarkUID(),
            Opacity: 0.8,
            Scale:   0.15,
        },
    }

    // Use TUS resumable upload for large files
    if video.FileSize > 100*1024*1024 { // 100MB
        return cs.resumableUpload(createReq, video)
    }

    return cs.directUpload(createReq, video)
}

func (cs *CloudflareStreamService) resumableUpload(req CloudflareCreateVideoRequest, video VideoUploadRequest) (*UploadResponse, error) {
    // Initialize TUS upload
    tusClient := tus.NewClient("https://upload.videodelivery.net/tus", &tus.Config{
        ChunkSize: 5 * 1024 * 1024, // 5MB chunks
        Header: map[string][]string{
            "Authorization": {fmt.Sprintf("Bearer %s", cs.apiToken)},
        },
    })

    upload := tus.NewUpload(video.File, tus.Config{
        Metadata: map[string]string{
            "filename":        video.Filename,
            "filetype":        video.ContentType,
            "requiresignedurls": "true",
        },
    })

    uploader := tus.NewUploader(tusClient, upload)
    return uploader.Upload()
}
```

#### Adaptive Streaming Configuration
```json
{
  "cloudflare_stream_settings": {
    "video_encoding": {
      "profiles": [
        {
          "name": "mobile_optimized",
          "resolutions": ["360p", "480p", "720p"],
          "bitrates": [800, 1200, 2500],
          "codecs": ["h264", "av1"],
          "target_devices": ["mobile", "tablet"]
        },
        {
          "name": "desktop_optimized",
          "resolutions": ["720p", "1080p", "1440p"],
          "bitrates": [2500, 5000, 8000],
          "codecs": ["h264", "hevc", "av1"],
          "target_devices": ["desktop", "smart_tv"]
        }
      ]
    },
    "adaptive_streaming": {
      "segment_duration": 6,
      "playlist_type": "hls",
      "enable_fast_start": true,
      "enable_low_latency": true
    }
  }
}
```

### 2. Client-Side Video Player Optimization

#### Optimized HLS.js Configuration
```javascript
// Enhanced video player with network intelligence
class IntelligentVideoPlayer {
    constructor(videoElement, streamUrl, sessionId) {
        this.video = videoElement;
        this.streamUrl = streamUrl;
        this.sessionId = sessionId;
        this.networkMonitor = new NetworkMonitor();
        this.websocket = new WebSocket(`wss://api.studyplatform.com/api/videos/ws/${sessionId}`);

        this.initializePlayer();
        this.setupNetworkMonitoring();
    }

    initializePlayer() {
        this.hls = new Hls({
            // Optimized configuration for adaptive streaming
            enableWorker: true,
            lowLatencyMode: true,
            backBufferLength: 30,
            maxBufferLength: 60,
            maxMaxBufferLength: 120,

            // Network intelligence integration
            abrEwmaFastLive: 3.0,
            abrEwmaSlowLive: 9.0,
            abrMaxWithRealBitrate: true,

            // Startup optimization
            startLevel: -1, // Auto-select based on bandwidth
            testBandwidth: true,

            // Error recovery
            manifestLoadingTimeOut: 20000,
            manifestLoadingMaxRetry: 6,
            levelLoadingTimeOut: 20000,
            levelLoadingMaxRetry: 6,

            // Fragment loading optimization
            fragLoadingTimeOut: 30000,
            fragLoadingMaxRetry: 6,
            fragLoadingRetryDelay: 1000,
        });

        this.hls.loadSource(this.streamUrl);
        this.hls.attachMedia(this.video);

        this.setupHLSEventHandlers();
    }

    setupNetworkMonitoring() {
        // Monitor network conditions every 5 seconds
        setInterval(() => {
            this.measureNetworkMetrics();
        }, 5000);

        // WebSocket message handling
        this.websocket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleServerRecommendation(message);
        };
    }

    async measureNetworkMetrics() {
        const metrics = await this.networkMonitor.measure();

        // Calculate buffer health
        const buffered = this.video.buffered;
        const currentTime = this.video.currentTime;
        let bufferHealth = 0;

        for (let i = 0; i < buffered.length; i++) {
            if (buffered.start(i) <= currentTime && currentTime <= buffered.end(i)) {
                bufferHealth = buffered.end(i) - currentTime;
                break;
            }
        }

        const networkData = {
            type: 'network_status',
            data: {
                bandwidth_mbps: metrics.bandwidth,
                latency_ms: metrics.latency,
                buffer_health_seconds: bufferHealth,
                current_quality: this.getCurrentQuality(),
                video_timestamp_seconds: this.video.currentTime,
                player_state: this.video.paused ? 'paused' : 'playing'
            }
        };

        this.websocket.send(JSON.stringify(networkData));
    }

    handleServerRecommendation(message) {
        switch (message.type) {
            case 'quality_recommendation':
                this.adaptQuality(message.data);
                break;
            case 'preload_instruction':
                this.preloadSegments(message.data);
                break;
            case 'buffer_optimization':
                this.optimizeBuffer(message.data);
                break;
        }
    }

    adaptQuality(recommendation) {
        const targetLevel = this.findLevelByLabel(recommendation.recommended_quality);

        if (targetLevel !== -1 && recommendation.should_apply_immediately) {
            // Smooth quality transition
            this.hls.nextLevel = targetLevel;
            this.hls.loadLevel = targetLevel;

            // Log quality change
            console.log(`Quality adapted to ${recommendation.recommended_quality} (${recommendation.reason})`);
        }
    }

    preloadSegments(instruction) {
        // Implement intelligent preloading based on server recommendations
        instruction.segments.forEach(segment => {
            const link = document.createElement('link');
            link.rel = 'prefetch';
            link.href = segment.url;
            document.head.appendChild(link);
        });
    }
}

// Network monitoring utility
class NetworkMonitor {
    async measure() {
        const startTime = performance.now();

        // Measure bandwidth using a small test file
        const response = await fetch('/api/bandwidth-test', {
            cache: 'no-cache'
        });

        const endTime = performance.now();
        const latency = endTime - startTime;

        // Estimate bandwidth based on response time and file size
        const contentLength = response.headers.get('content-length');
        const bandwidth = contentLength ? (contentLength * 8) / (latency * 1000) : 0;

        return {
            bandwidth: bandwidth,
            latency: latency,
            timestamp: Date.now()
        };
    }
}
```

## Scalability and Performance Patterns

### 1. Microservice Performance Optimization

#### Circuit Breaker Implementation
```go
type CircuitBreaker struct {
    maxRequests  uint32
    interval     time.Duration
    timeout      time.Duration
    readyToTrip  func(counts Counts) bool
    onStateChange func(name string, from State, to State)

    mutex   sync.Mutex
    state   State
    counts  Counts
    expiry  time.Time
}

func (cb *CircuitBreaker) Call(req func() (interface{}, error)) (interface{}, error) {
    generation, err := cb.beforeRequest()
    if err != nil {
        return nil, err
    }

    defer func() {
        e := recover()
        if e != nil {
            cb.afterRequest(generation, false)
            panic(e)
        }
    }()

    result, err := req()
    cb.afterRequest(generation, err == nil)
    return result, err
}

// Video service with circuit breakers
type VideoServiceHandler struct {
    authBreaker      *CircuitBreaker
    cloudflareBreaker *CircuitBreaker
    redisBreaker     *CircuitBreaker
}

func (h *VideoServiceHandler) GetVideoWithAuth(ctx context.Context, videoID, userID uuid.UUID) (*Video, error) {
    // Use circuit breaker for auth validation
    authResult, err := h.authBreaker.Call(func() (interface{}, error) {
        return h.authService.ValidateVideoAccess(ctx, userID, videoID)
    })

    if err != nil {
        // Fallback to cached permissions or basic validation
        return h.getVideoWithFallbackAuth(videoID, userID)
    }

    return h.getVideoWithPermissions(videoID, authResult.(*AuthResult))
}
```

#### Connection Pool Optimization
```go
type ServiceConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    HTTP     HTTPConfig     `yaml:"http"`
}

type DatabaseConfig struct {
    MaxOpenConns    int           `yaml:"max_open_conns"`
    MaxIdleConns    int           `yaml:"max_idle_conns"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
    ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

type RedisConfig struct {
    PoolSize     int           `yaml:"pool_size"`
    MinIdleConns int           `yaml:"min_idle_conns"`
    MaxConnAge   time.Duration `yaml:"max_conn_age"`
    PoolTimeout  time.Duration `yaml:"pool_timeout"`
}

// Optimized connection pool setup
func setupDatabasePool(config DatabaseConfig) *sql.DB {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    db.SetMaxOpenConns(config.MaxOpenConns)     // 100
    db.SetMaxIdleConns(config.MaxIdleConns)     // 25
    db.SetConnMaxLifetime(config.ConnMaxLifetime) // 1 hour
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 15 minutes

    return db
}

func setupRedisPool(config RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        PoolSize:     config.PoolSize,     // 50
        MinIdleConns: config.MinIdleConns, // 10
        MaxConnAge:   config.MaxConnAge,   // 30 minutes
        PoolTimeout:  config.PoolTimeout,  // 4 seconds
    })
}
```

### 2. Caching Strategy

#### Multi-Level Caching
```go
type VideoCache struct {
    l1Cache    *ristretto.Cache // In-memory cache
    l2Cache    *redis.Client    // Redis cache
    database   *sql.DB         // Database
}

func (vc *VideoCache) GetVideo(videoID uuid.UUID) (*Video, error) {
    // L1 Cache (in-memory)
    if video, found := vc.l1Cache.Get(videoID.String()); found {
        return video.(*Video), nil
    }

    // L2 Cache (Redis)
    videoJSON, err := vc.l2Cache.Get(context.Background(),
        fmt.Sprintf("video:%s", videoID)).Result()
    if err == nil {
        var video Video
        if json.Unmarshal([]byte(videoJSON), &video) == nil {
            // Store in L1 cache
            vc.l1Cache.Set(videoID.String(), &video, time.Hour)
            return &video, nil
        }
    }

    // Database
    video, err := vc.getVideoFromDatabase(videoID)
    if err != nil {
        return nil, err
    }

    // Cache in both levels
    go func() {
        videoJSON, _ := json.Marshal(video)
        vc.l2Cache.Set(context.Background(),
            fmt.Sprintf("video:%s", videoID), videoJSON, time.Hour*6)
        vc.l1Cache.Set(videoID.String(), video, time.Hour)
    }()

    return video, nil
}
```

### 3. Load Balancing and Auto-scaling

#### Horizontal Pod Autoscaler Configuration
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: video-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: video-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: concurrent_video_sessions
      target:
        type: AverageValue
        averageValue: "100"
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
```

## Performance Monitoring and Alerting

### 1. Key Performance Metrics

```yaml
# Prometheus metrics configuration
video_streaming_metrics:
  - name: video_startup_time_seconds
    type: histogram
    buckets: [0.5, 1, 2, 3, 5, 10]
    labels: [quality, device_type, connection_type]

  - name: video_buffering_events_total
    type: counter
    labels: [video_id, quality, reason]

  - name: network_quality_score
    type: gauge
    labels: [user_id, connection_type]

  - name: concurrent_video_sessions
    type: gauge
    labels: [video_id, quality]

  - name: quality_changes_total
    type: counter
    labels: [from_quality, to_quality, reason, auto_switch]
```

### 2. Performance Alerting Rules

```yaml
# Alert rules for video streaming performance
groups:
- name: video_streaming_alerts
  rules:
  - alert: HighVideoStartupTime
    expr: histogram_quantile(0.95, video_startup_time_seconds) > 3
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High video startup time detected"
      description: "95th percentile startup time is {{ $value }}s"

  - alert: HighBufferingRate
    expr: rate(video_buffering_events_total[5m]) > 0.1
    for: 3m
    labels:
      severity: critical
    annotations:
      summary: "High buffering event rate"
      description: "Buffering rate is {{ $value }} events per second"

  - alert: LowNetworkQualityScore
    expr: avg(network_quality_score) < 5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Low average network quality score"
      description: "Average quality score is {{ $value }}"
```

This comprehensive performance optimization guide ensures the video streaming platform can handle high loads while maintaining excellent user experience with sub-3-second startup times and minimal buffering events.
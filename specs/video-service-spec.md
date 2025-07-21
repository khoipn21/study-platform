# Video Service Specification

## Overview
The Video Service manages video streaming, processing, and delivery using Cloudflare Stream for video hosting and Redis for message queuing to track user network status and optimize streaming quality. It provides adaptive bitrate streaming, analytics, and intelligent quality adjustment based on user network conditions.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌──────────────────┐
│   API Gateway   │───▶│  Video Service  │───▶│ Cloudflare Stream│
└─────────────────┘    └─────────────────┘    └──────────────────┘
                              │
                              ├─────────────────▶┌─────────────────┐
                              │                  │     Redis       │
                              │                  │ (Message Queue) │
                              │                  └─────────────────┘
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │  (Video Meta)   │
                       └─────────────────┘
```

## Technology Stack
- **Language**: Go 1.21+
- **Video CDN**: Cloudflare Stream
- **Message Queue**: Redis (for network status tracking)
- **Database**: PostgreSQL (video metadata)
- **WebSocket**: For real-time network monitoring
- **Authentication**: JWT tokens from Auth Service
- **Container**: Docker

## Features

### Core Features
1. **Video Management**
   - Video upload and processing
   - Metadata management
   - Video transcoding via Cloudflare
   - Thumbnail generation

2. **Adaptive Streaming**
   - Multiple bitrate options
   - Quality auto-adjustment
   - Network-aware streaming
   - Buffer optimization

3. **Network Intelligence**
   - Real-time bandwidth detection
   - Connection quality monitoring
   - Smart quality switching
   - Preloading optimization

4. **Analytics & Monitoring**
   - Watch time tracking
   - Quality metrics
   - User engagement analytics
   - Performance monitoring

5. **Access Control**
   - Video permissions
   - Token-based access
   - Time-limited URLs
   - Geographic restrictions

## Database Schema

```sql
-- Video metadata table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cloudflare_uid VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_seconds INTEGER,
    file_size_bytes BIGINT,
    upload_user_id UUID NOT NULL,
    course_id UUID,
    lecture_id UUID,
    status VARCHAR(20) DEFAULT 'processing', -- 'uploading', 'processing', 'ready', 'error'
    visibility VARCHAR(20) DEFAULT 'private', -- 'public', 'private', 'unlisted'
    thumbnail_url VARCHAR(500),
    stream_url VARCHAR(500),
    preview_url VARCHAR(500),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Video quality variants
CREATE TABLE video_qualities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    quality_label VARCHAR(10) NOT NULL, -- '360p', '720p', '1080p'
    bitrate_kbps INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    fps INTEGER DEFAULT 30,
    codec VARCHAR(20) DEFAULT 'h264',
    url VARCHAR(500) NOT NULL,
    file_size_bytes BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User viewing sessions
CREATE TABLE viewing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    started_at TIMESTAMP DEFAULT NOW(),
    last_heartbeat TIMESTAMP DEFAULT NOW(),
    current_time_seconds INTEGER DEFAULT 0,
    current_quality VARCHAR(10),
    total_watch_time_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    ip_address INET,
    created_at TIMESTAMP DEFAULT NOW()
);

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

-- Video access permissions
CREATE TABLE video_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID,
    role_id UUID,
    permission_type VARCHAR(20) NOT NULL, -- 'view', 'download', 'share'
    granted_by UUID NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Analytics aggregations
CREATE TABLE video_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_views INTEGER DEFAULT 0,
    unique_viewers INTEGER DEFAULT 0,
    total_watch_time_seconds BIGINT DEFAULT 0,
    avg_watch_time_seconds INTEGER DEFAULT 0,
    completion_rate DECIMAL(5,2) DEFAULT 0,
    quality_distribution JSONB, -- {'360p': 30, '720p': 50, '1080p': 20}
    geographic_distribution JSONB,
    device_distribution JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(video_id, date)
);

-- Indexes
CREATE INDEX idx_videos_cloudflare_uid ON videos(cloudflare_uid);
CREATE INDEX idx_videos_user_id ON videos(upload_user_id);
CREATE INDEX idx_videos_course_lecture ON videos(course_id, lecture_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_video_qualities_video_id ON video_qualities(video_id);
CREATE INDEX idx_viewing_sessions_user_video ON viewing_sessions(user_id, video_id);
CREATE INDEX idx_viewing_sessions_session ON viewing_sessions(session_id);
CREATE INDEX idx_network_metrics_session ON network_metrics(session_id);
CREATE INDEX idx_network_metrics_timestamp ON network_metrics(timestamp);
CREATE INDEX idx_video_permissions_video_user ON video_permissions(video_id, user_id);
CREATE INDEX idx_video_analytics_video_date ON video_analytics(video_id, date);
```

## Redis Message Queue Schema

### Channel Patterns
```redis
# Network status updates
channel: network_status:{session_id}
payload: {
    "session_id": "string",
    "user_id": "uuid",
    "timestamp": "2024-01-01T12:00:00Z",
    "bandwidth_mbps": 5.2,
    "latency_ms": 150,
    "packet_loss": 0.1,
    "connection_type": "wifi",
    "quality_score": 7,
    "current_quality": "720p",
    "recommended_quality": "720p",
    "buffer_health": 8
}

# Quality change requests
channel: quality_change:{session_id}
payload: {
    "session_id": "string",
    "video_id": "uuid",
    "from_quality": "720p",
    "to_quality": "480p",
    "reason": "bandwidth_drop",
    "timestamp": "2024-01-01T12:00:00Z"
}

# Viewing analytics
channel: video_analytics
payload: {
    "event_type": "watch_time",
    "session_id": "string",
    "video_id": "uuid",
    "user_id": "uuid",
    "current_time": 120,
    "quality": "720p",
    "timestamp": "2024-01-01T12:00:00Z"
}

# Heartbeat for active sessions
channel: session_heartbeat
payload: {
    "session_id": "string",
    "user_id": "uuid",
    "video_id": "uuid",
    "current_time": 300,
    "quality": "720p",
    "buffer_health": 5,
    "timestamp": "2024-01-01T12:00:00Z"
}
```

## API Endpoints

### Video Upload
```http
POST /api/videos/upload
Content-Type: multipart/form-data

Form fields:
- video: binary video file
- title: string
- description: string (optional)
- course_id: uuid (optional)
- lecture_id: uuid (optional)
- visibility: string (public/private/unlisted)

Response:
{
    "video_id": "uuid",
    "cloudflare_uid": "string",
    "title": "Video Title",
    "status": "processing",
    "upload_url": "https://upload.videodelivery.net/...",
    "estimated_processing_time": "5-10 minutes"
}
```

### Get Video Details
```http
GET /api/videos/{video_id}

Response:
{
    "id": "uuid",
    "cloudflare_uid": "string",
    "title": "Video Title",
    "description": "Video description",
    "duration": 1800,
    "status": "ready",
    "visibility": "private",
    "thumbnail_url": "https://cloudflare.../thumb.jpg",
    "qualities": [
        {
            "label": "360p",
            "bitrate": 800,
            "width": 640,
            "height": 360,
            "url": "https://cloudflare.../360p.m3u8"
        }
    ],
    "created_at": "2024-01-01T12:00:00Z"
}
```

### Start Viewing Session
```http
POST /api/videos/{video_id}/sessions
Content-Type: application/json

{
    "device_info": {
        "user_agent": "string",
        "screen_resolution": "1920x1080",
        "connection_type": "wifi"
    }
}

Response:
{
    "session_id": "uuid",
    "stream_url": "https://cloudflare.../stream.m3u8",
    "thumbnail_url": "https://cloudflare.../thumb.jpg",
    "qualities": [...],
    "recommended_quality": "720p",
    "websocket_url": "wss://video-service/ws/{session_id}",
    "expires_at": "2024-01-01T18:00:00Z"
}
```

### Update Network Status
```http
POST /api/videos/sessions/{session_id}/network
Content-Type: application/json

{
    "bandwidth_mbps": 5.2,
    "latency_ms": 150,
    "packet_loss": 0.1,
    "connection_type": "wifi",
    "buffer_health": 8,
    "current_time": 120,
    "current_quality": "720p"
}

Response:
{
    "recommended_quality": "720p",
    "quality_score": 7,
    "should_preload": true,
    "buffer_target": 10
}
```

### Video Analytics
```http
GET /api/videos/{video_id}/analytics?period=7d

Response:
{
    "video_id": "uuid",
    "period": "7d",
    "total_views": 1250,
    "unique_viewers": 890,
    "total_watch_time": 156000,
    "avg_watch_time": 175,
    "completion_rate": 0.65,
    "daily_stats": [
        {
            "date": "2024-01-01",
            "views": 180,
            "unique_viewers": 150,
            "watch_time": 22500
        }
    ],
    "quality_distribution": {
        "360p": 25,
        "720p": 55,
        "1080p": 20
    }
}
```

## WebSocket Protocol

### Connection
```
ws://video-service/ws/{session_id}?token={jwt_token}
```

### Message Types

#### Client to Server
```json
// Network quality update
{
    "type": "network_status",
    "data": {
        "bandwidth_mbps": 5.2,
        "latency_ms": 150,
        "packet_loss": 0.1,
        "buffer_health": 8,
        "current_time": 120,
        "current_quality": "720p"
    }
}

// Viewing progress update
{
    "type": "progress_update",
    "data": {
        "current_time": 300,
        "quality": "720p",
        "paused": false
    }
}

// Quality change request
{
    "type": "quality_change",
    "data": {
        "requested_quality": "480p",
        "reason": "user_preference"
    }
}
```

#### Server to Client
```json
// Quality recommendation
{
    "type": "quality_recommendation",
    "data": {
        "recommended_quality": "480p",
        "reason": "bandwidth_drop",
        "confidence": 0.85
    }
}

// Preload instruction
{
    "type": "preload",
    "data": {
        "segments": ["seg_120.ts", "seg_121.ts"],
        "priority": "high"
    }
}

// Analytics trigger
{
    "type": "analytics_event",
    "data": {
        "event": "quality_change",
        "from": "720p",
        "to": "480p"
    }
}
```

## Configuration

### Environment Variables
```env
# Service Configuration
VIDEO_SERVICE_PORT=8085
VIDEO_SERVICE_HOST=0.0.0.0

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=video_user
DB_PASSWORD=video_password
DB_NAME=video_service

# Cloudflare Stream Configuration
CLOUDFLARE_ACCOUNT_ID=your-account-id
CLOUDFLARE_STREAM_TOKEN=your-stream-token
CLOUDFLARE_API_EMAIL=your-email@domain.com
CLOUDFLARE_API_KEY=your-api-key

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
REDIS_DB=0
REDIS_POOL_SIZE=10

# WebSocket Configuration
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
WS_HEARTBEAT_INTERVAL=30s

# Video Processing
MAX_VIDEO_SIZE=5368709120  # 5GB
ALLOWED_FORMATS=mp4,avi,mov,wmv,flv,webm
DEFAULT_THUMBNAIL_TIME=10  # seconds

# Network Intelligence
BANDWIDTH_CHECK_INTERVAL=10s
QUALITY_CHANGE_THRESHOLD=0.8
BUFFER_TARGET_SECONDS=10
MIN_BUFFER_SECONDS=5

# Security
JWT_SECRET=your-jwt-secret
MAX_SESSION_DURATION=24h
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com

# Analytics
ANALYTICS_BATCH_SIZE=100
ANALYTICS_FLUSH_INTERVAL=60s
```

## Service Implementation

### Directory Structure
```
video-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── video.go
│   │   ├── session.go
│   │   ├── analytics.go
│   │   └── websocket.go
│   ├── service/
│   │   ├── video.go
│   │   ├── cloudflare.go
│   │   ├── network_intelligence.go
│   │   ├── analytics.go
│   │   └── streaming.go
│   ├── repository/
│   │   ├── video_repository.go
│   │   ├── session_repository.go
│   │   └── analytics_repository.go
│   ├── queue/
│   │   ├── redis_client.go
│   │   ├── publisher.go
│   │   └── subscriber.go
│   ├── websocket/
│   │   ├── hub.go
│   │   ├── client.go
│   │   └── message.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logging.go
│   └── model/
│       ├── video.go
│       ├── session.go
│       └── analytics.go
├── migrations/
│   ├── 001_create_videos_table.up.sql
│   ├── 002_create_video_qualities_table.up.sql
│   ├── 003_create_viewing_sessions_table.up.sql
│   ├── 004_create_network_metrics_table.up.sql
│   └── 005_create_video_analytics_table.up.sql
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Implementation Steps

### Phase 1: Core Setup
1. **Project Structure**
   ```bash
   mkdir -p video-service/{cmd/server,internal/{config,handler,service,repository,queue,websocket,middleware,model},migrations}
   cd video-service && go mod init video-service
   ```

2. **Dependencies**
   ```bash
   go get github.com/cloudflare/cloudflare-go
   go get github.com/gin-gonic/gin
   go get github.com/gorilla/websocket
   go get github.com/go-redis/redis/v8
   go get github.com/lib/pq
   go get gorm.io/gorm
   go get github.com/golang-migrate/migrate/v4
   ```

3. **Database & Redis Setup**
   - Create migration files
   - Set up GORM models
   - Configure Redis client

### Phase 2: Cloudflare Integration
1. **Stream API Integration**
   - Upload video to Cloudflare
   - Webhook handling for processing status
   - Thumbnail generation

2. **Video Management**
   - CRUD operations for videos
   - Quality variant management
   - Access control

### Phase 3: Network Intelligence
1. **WebSocket Implementation**
   - Real-time connection handling
   - Message routing
   - Connection management

2. **Network Monitoring**
   - Bandwidth detection
   - Quality scoring algorithm
   - Adaptive streaming logic

### Phase 4: Message Queue System
1. **Redis Integration**
   - Publisher/subscriber pattern
   - Message serialization
   - Queue management

2. **Analytics Processing**
   - Real-time data collection
   - Batch processing
   - Data aggregation

### Phase 5: Advanced Features
1. **Smart Streaming**
   - Preloading optimization
   - Buffer management
   - Quality switching

2. **Analytics Dashboard**
   - Real-time metrics
   - Historical data
   - Performance insights

## Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o video-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/video-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8085
CMD ["./video-service"]
```

### Docker Compose Service
```yaml
video-service:
  build:
    context: ./video-service
    dockerfile: Dockerfile
  ports:
    - "8085:8085"
  environment:
    - DB_HOST=postgres
    - DB_PORT=5432
    - DB_USER=video_user
    - DB_PASSWORD=video_password
    - DB_NAME=video_service
    - REDIS_HOST=redis
    - REDIS_PORT=6379
    - CLOUDFLARE_ACCOUNT_ID=${CLOUDFLARE_ACCOUNT_ID}
    - CLOUDFLARE_STREAM_TOKEN=${CLOUDFLARE_STREAM_TOKEN}
  depends_on:
    - postgres
    - redis
  volumes:
    - ./video-service/migrations:/app/migrations
  networks:
    - study-platform

redis:
  image: redis:7-alpine
  ports:
    - "6379:6379"
  command: redis-server --requirepass redis_password
  volumes:
    - redis_data:/data
  networks:
    - study-platform

volumes:
  redis_data:
```

## Network Intelligence Algorithm

### Quality Scoring
```go
func CalculateQualityScore(metrics NetworkMetrics) int {
    score := 10
    
    // Bandwidth factor (40% weight)
    if metrics.BandwidthMbps < 1.0 {
        score -= 4
    } else if metrics.BandwidthMbps < 3.0 {
        score -= 2
    } else if metrics.BandwidthMbps < 5.0 {
        score -= 1
    }
    
    // Latency factor (30% weight)
    if metrics.LatencyMs > 500 {
        score -= 3
    } else if metrics.LatencyMs > 200 {
        score -= 2
    } else if metrics.LatencyMs > 100 {
        score -= 1
    }
    
    // Packet loss factor (20% weight)
    if metrics.PacketLoss > 0.05 {
        score -= 2
    } else if metrics.PacketLoss > 0.01 {
        score -= 1
    }
    
    // Buffer health factor (10% weight)
    if metrics.BufferHealth < 3 {
        score -= 1
    }
    
    if score < 1 {
        score = 1
    }
    
    return score
}

func RecommendQuality(score int, currentQuality string) string {
    qualityMap := map[int]string{
        1: "240p", 2: "240p", 3: "360p",
        4: "360p", 5: "480p", 6: "480p",
        7: "720p", 8: "720p", 9: "1080p", 10: "1080p",
    }
    
    recommended := qualityMap[score]
    
    // Avoid frequent switching
    if shouldPreventSwitch(currentQuality, recommended) {
        return currentQuality
    }
    
    return recommended
}
```

## Performance Optimization

### Caching Strategy
- Redis caching for video metadata
- CDN caching for thumbnails
- Quality recommendation caching
- Session data caching

### Database Optimization
- Proper indexing for queries
- Read replicas for analytics
- Connection pooling
- Query optimization

### Network Optimization
- WebSocket connection pooling
- Message batching
- Compression for large payloads
- Heartbeat optimization

## Monitoring & Alerting

### Key Metrics
- Video processing success rate
- Streaming quality distribution
- User engagement metrics
- Network performance scores
- WebSocket connection health

### Alert Conditions
- High video processing failures
- Poor network quality scores
- Excessive quality switching
- WebSocket connection drops
- Redis queue backup

## Security Considerations

### Video Security
- Signed URLs with expiration
- Domain restrictions
- Geographic blocking
- Token validation

### WebSocket Security
- JWT token validation
- Rate limiting per connection
- Message validation
- Connection throttling

### Data Privacy
- User data anonymization
- GDPR compliance
- Data retention policies
- Audit logging

## Deployment Checklist

- [ ] Configure Cloudflare Stream account
- [ ] Set up Redis cluster
- [ ] Configure database with proper indexes
- [ ] Set up monitoring and alerting
- [ ] Configure WebSocket load balancing
- [ ] Test video upload and streaming
- [ ] Performance benchmarking
- [ ] Security audit
- [ ] Disaster recovery testing
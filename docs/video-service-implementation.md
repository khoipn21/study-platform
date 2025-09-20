# Video Service Technical Implementation

## Service Architecture Overview

The enhanced video streaming platform uses a microservices architecture with specialized services for handling different aspects of video streaming, network intelligence, and real-time adaptation.

### Service Interaction Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │  Auth Service   │    │ Course Service  │
│    (Port 8080)  │    │   (Port 8081)   │    │  (Port 8082)    │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP/WebSocket       │ gRPC                 │ gRPC
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Video Service (Port 8084)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Video     │  │   Session   │  │    Network Intelligence │ │
│  │ Management  │  │  Tracking   │  │        Engine           │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────┬───────────────────┬─────────────────────┬───────────┘
          │                   │                     │
          │ HTTP              │ WebSocket           │ Redis Streams
          │                   │                     │
          ▼                   ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Cloudflare      │  │   WebSocket     │  │     Redis       │
│   Stream API    │  │     Hub         │  │   Streams       │
│                 │  │                 │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
          │                   │                     │
          │                   │ Real-time           │ Event
          │ Webhooks         │ Updates             │ Processing
          │                   │                     │
          ▼                   ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   PostgreSQL    │  │     Client      │  │   Analytics     │
│   Database      │  │  Video Player   │  │   Service       │
│                 │  │                 │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## Enhanced API Specifications

### 1. Video Management APIs

#### POST /api/videos/upload
**Purpose**: Initialize video upload to Cloudflare Stream

**Request Body**:
```json
{
  "title": "Introduction to Microservices",
  "description": "Learn the fundamentals of microservices architecture",
  "course_id": "course_uuid",
  "lecture_id": "lecture_uuid",
  "visibility": "course_only",
  "metadata": {
    "instructor_id": "user_uuid",
    "tags": ["programming", "architecture"],
    "expected_duration": 1800
  }
}
```

**Response**:
```json
{
  "video_id": "video_uuid",
  "cloudflare_uid": "cf_video_id",
  "upload_url": "https://upload.videodelivery.net/...",
  "upload_token": "upload_token",
  "estimated_processing_time": "5-10 minutes",
  "webhook_url": "https://api.studyplatform.com/api/videos/webhooks/cloudflare",
  "max_file_size_bytes": 10737418240,
  "supported_formats": ["mp4", "mov", "avi", "mkv"]
}
```

#### GET /api/videos/{video_id}
**Purpose**: Get video metadata and streaming information

**Response**:
```json
{
  "id": "video_uuid",
  "cloudflare_uid": "cf_video_id",
  "title": "Introduction to Microservices",
  "description": "Learn the fundamentals...",
  "status": "ready",
  "visibility": "course_only",
  "duration_seconds": 1847,
  "file_size_bytes": 524288000,
  "thumbnail_url": "https://videodelivery.net/thumb.jpg",
  "stream_url": "https://videodelivery.net/manifest.m3u8",
  "preview_url": "https://videodelivery.net/preview.mp4",
  "qualities": [
    {
      "quality_label": "360p",
      "bitrate_kbps": 800,
      "width": 640,
      "height": 360,
      "url": "https://videodelivery.net/360p.m3u8"
    },
    {
      "quality_label": "720p",
      "bitrate_kbps": 2500,
      "width": 1280,
      "height": 720,
      "url": "https://videodelivery.net/720p.m3u8"
    },
    {
      "quality_label": "1080p",
      "bitrate_kbps": 5000,
      "width": 1920,
      "height": 1080,
      "url": "https://videodelivery.net/1080p.m3u8"
    }
  ],
  "metadata": {
    "instructor_id": "user_uuid",
    "tags": ["programming", "architecture"],
    "chapters": [
      {"title": "Introduction", "start_time": 0},
      {"title": "Core Concepts", "start_time": 300},
      {"title": "Implementation", "start_time": 900}
    ]
  },
  "created_at": "2024-09-17T10:00:00Z",
  "updated_at": "2024-09-17T10:15:00Z"
}
```

### 2. Video Session Management APIs

#### POST /api/videos/{video_id}/sessions
**Purpose**: Create a new viewing session with network intelligence

**Request Body**:
```json
{
  "device_info": {
    "user_agent": "Mozilla/5.0...",
    "screen_resolution": "1920x1080",
    "connection_type": "wifi",
    "device_type": "desktop",
    "supports_hls": true,
    "supports_dash": false
  },
  "initial_quality": "auto"
}
```

**Response**:
```json
{
  "session_id": "sess_unique_id",
  "stream_url": "https://videodelivery.net/signed_manifest.m3u8",
  "thumbnail_url": "https://videodelivery.net/thumb.jpg",
  "qualities": [...],
  "recommended_quality": "720p",
  "websocket_url": "wss://api.studyplatform.com/api/videos/ws/sess_unique_id",
  "expires_at": "2024-09-17T18:00:00Z",
  "buffer_settings": {
    "initial_buffer_seconds": 10,
    "max_buffer_seconds": 30,
    "rebuffer_threshold_seconds": 3
  }
}
```

#### PUT /api/videos/sessions/{session_id}/progress
**Purpose**: Update viewing progress and collect metrics

**Request Body**:
```json
{
  "current_time_seconds": 450,
  "current_quality": "720p",
  "network_metrics": {
    "bandwidth_mbps": 12.5,
    "latency_ms": 55,
    "packet_loss_percent": 0.2,
    "buffer_health_seconds": 15,
    "jitter_ms": 8
  },
  "player_state": "playing",
  "buffering_events": [
    {
      "started_at": "2024-09-17T10:30:00Z",
      "duration_ms": 2500,
      "video_timestamp": 445
    }
  ]
}
```

**Response**:
```json
{
  "quality_recommendation": {
    "recommended_quality": "720p",
    "confidence": 0.87,
    "reason": "stable_connection"
  },
  "buffer_strategy": {
    "should_preload": true,
    "preload_seconds": 20,
    "target_buffer_seconds": 15
  },
  "session_updated": true
}
```

### 3. Network Intelligence APIs

#### POST /api/videos/sessions/{session_id}/network
**Purpose**: Report detailed network metrics for adaptive streaming

**Request Body**:
```json
{
  "metrics": [
    {
      "timestamp": "2024-09-17T10:30:00Z",
      "bandwidth_mbps": 15.2,
      "latency_ms": 45,
      "packet_loss_percent": 0.1,
      "jitter_ms": 5,
      "connection_type": "wifi",
      "signal_strength": -45,
      "throughput_mbps": 14.8
    }
  ],
  "quality_events": [
    {
      "timestamp": "2024-09-17T10:30:05Z",
      "from_quality": "1080p",
      "to_quality": "720p",
      "reason": "bandwidth_drop",
      "auto_switch": true
    }
  ]
}
```

**Response**:
```json
{
  "quality_score": 7,
  "recommendations": {
    "immediate": {
      "quality": "720p",
      "reason": "bandwidth_optimization",
      "confidence": 0.92
    },
    "predictive": {
      "quality": "1080p",
      "condition": "if_bandwidth_stable_for_30s",
      "confidence": 0.75
    }
  },
  "buffer_guidance": {
    "target_seconds": 15,
    "preload_next_segments": 3
  }
}
```

### 4. Analytics and Monitoring APIs

#### GET /api/videos/{video_id}/analytics
**Purpose**: Get comprehensive video analytics

**Query Parameters**:
- `period`: `hour`, `day`, `week`, `month`
- `start_date`: ISO 8601 date
- `end_date`: ISO 8601 date
- `granularity`: `minute`, `hour`, `day`

**Response**:
```json
{
  "video_id": "video_uuid",
  "period": {
    "start": "2024-09-10T00:00:00Z",
    "end": "2024-09-17T23:59:59Z"
  },
  "summary": {
    "total_views": 1247,
    "unique_viewers": 892,
    "total_watch_time_seconds": 2456789,
    "avg_watch_time_seconds": 1847,
    "completion_rate": 67.5,
    "peak_concurrent_viewers": 234
  },
  "quality_metrics": {
    "avg_quality_score": 7.8,
    "quality_distribution": {
      "360p": 15.2,
      "720p": 45.8,
      "1080p": 39.0
    },
    "total_quality_changes": 3456,
    "auto_quality_changes": 2987
  },
  "performance_metrics": {
    "avg_startup_time_seconds": 2.3,
    "buffering_incident_rate": 8.5,
    "avg_buffering_duration_seconds": 3.2,
    "error_rate": 0.8
  },
  "engagement_metrics": {
    "play_events": 1247,
    "pause_events": 2345,
    "seek_events": 4567,
    "fullscreen_events": 678,
    "replay_rate": 12.3
  },
  "geographic_distribution": {
    "US": 52.3,
    "EU": 28.7,
    "AS": 15.2,
    "Other": 3.8
  },
  "device_distribution": {
    "desktop": 65.2,
    "mobile": 28.9,
    "tablet": 5.9
  },
  "time_series": [
    {
      "timestamp": "2024-09-17T10:00:00Z",
      "concurrent_viewers": 156,
      "avg_quality_score": 7.9,
      "buffering_incidents": 12
    }
  ]
}
```

### 5. WebSocket Real-time Communication

#### Connection Establishment
```
WebSocket URL: wss://api.studyplatform.com/api/videos/ws/{session_id}
Protocol: video-streaming-v1
```

#### Message Types

**Client → Server: Network Status Update**
```json
{
  "type": "network_status",
  "timestamp": "2024-09-17T10:30:00Z",
  "data": {
    "bandwidth_mbps": 15.2,
    "latency_ms": 45,
    "buffer_health_seconds": 12,
    "current_quality": "1080p",
    "player_state": "playing",
    "video_timestamp_seconds": 450
  }
}
```

**Server → Client: Quality Recommendation**
```json
{
  "type": "quality_recommendation",
  "timestamp": "2024-09-17T10:30:01Z",
  "data": {
    "recommended_quality": "720p",
    "reason": "bandwidth_optimization",
    "confidence": 0.92,
    "should_apply_immediately": true,
    "estimated_improvement": "reduce_buffering_by_60%"
  }
}
```

**Server → Client: Preload Instructions**
```json
{
  "type": "preload_instruction",
  "timestamp": "2024-09-17T10:30:02Z",
  "data": {
    "segments": [
      {
        "url": "segment_123.m4s",
        "quality": "720p",
        "start_time": 460,
        "duration": 10
      }
    ],
    "priority": "high",
    "cache_duration_seconds": 300
  }
}
```

**Bidirectional: Error Notification**
```json
{
  "type": "error",
  "timestamp": "2024-09-17T10:30:03Z",
  "data": {
    "error_code": "NETWORK_ERROR",
    "error_message": "Connection timeout",
    "video_timestamp_seconds": 455,
    "suggested_action": "retry_with_lower_quality",
    "recovery_instructions": {
      "fallback_quality": "480p",
      "retry_delay_seconds": 5
    }
  }
}
```

### 6. Cloudflare Stream Webhook Handler

#### POST /api/videos/webhooks/cloudflare
**Purpose**: Handle video processing status updates from Cloudflare

**Request Body (from Cloudflare)**:
```json
{
  "uid": "cf_video_id",
  "status": {
    "state": "ready",
    "pctComplete": "100",
    "errorReasonCode": "",
    "errorReasonText": ""
  },
  "meta": {
    "name": "Introduction to Microservices"
  },
  "created": "2024-09-17T10:00:00Z",
  "modified": "2024-09-17T10:15:00Z",
  "size": 524288000,
  "duration": 1847,
  "input": {
    "width": 1920,
    "height": 1080
  },
  "playback": {
    "hls": "https://videodelivery.net/cf_video_id/manifest/video.m3u8",
    "dash": "https://videodelivery.net/cf_video_id/manifest/video.mpd"
  },
  "thumbnail": "https://videodelivery.net/cf_video_id/thumbnails/thumbnail.jpg",
  "preview": "https://videodelivery.net/cf_video_id/preview/preview.mp4"
}
```

## Database Repository Patterns

### 1. Video Repository Enhanced Methods

```go
type VideoRepository interface {
    // Basic CRUD
    CreateVideo(video *model.Video) error
    GetVideoByID(id uuid.UUID) (*model.Video, error)
    UpdateVideo(video *model.Video) error
    DeleteVideo(id uuid.UUID) error

    // Advanced queries
    GetVideoWithQualities(id uuid.UUID) (*model.VideoWithQualities, error)
    GetVideosByStatus(status string, limit, offset int) ([]*model.Video, error)
    SearchVideosWithFilters(filters VideoSearchFilters) ([]*model.Video, error)
    GetVideoAnalyticsSummary(videoID uuid.UUID, period time.Duration) (*model.VideoAnalytics, error)

    // Session management
    CreateViewingSession(session *model.ViewingSession) error
    UpdateSessionProgress(sessionID string, progress SessionProgress) error
    GetActiveSessionsByVideo(videoID uuid.UUID) ([]*model.ViewingSession, error)

    // Network metrics
    CreateNetworkMetrics(metrics *model.NetworkMetrics) error
    GetNetworkMetricsHistory(sessionID string, duration time.Duration) ([]*model.NetworkMetrics, error)
    GetAverageQualityScore(videoID uuid.UUID, period time.Duration) (float64, error)

    // Analytics
    RecordQualityChange(change *model.QualityChange) error
    RecordEngagementEvent(event *model.VideoEngagementEvent) error
    GetVideoPerformanceMetrics(videoID uuid.UUID) (*model.VideoPerformance, error)
}
```

### 2. Redis Stream Integration

```go
type RedisStreamService interface {
    // Network quality streams
    PublishNetworkMetrics(userID string, metrics NetworkMetrics) error
    ConsumeNetworkMetrics(consumerGroup string, handler NetworkMetricsHandler) error

    // Quality adaptation streams
    PublishQualityChange(change QualityChangeEvent) error
    ConsumeQualityChanges(consumerGroup string, handler QualityChangeHandler) error

    // Real-time analytics
    PublishAnalyticsEvent(event AnalyticsEvent) error
    ConsumeAnalyticsEvents(consumerGroup string, handler AnalyticsHandler) error

    // Stream management
    CreateConsumerGroup(stream, group string) error
    GetStreamInfo(stream string) (*StreamInfo, error)
    TrimStream(stream string, maxLength int64) error
}
```

## Service Communication Patterns

### 1. gRPC Integration with Auth Service

```protobuf
service VideoAuthService {
    rpc ValidateVideoAccess(VideoAccessRequest) returns (VideoAccessResponse);
    rpc GetUserPermissions(UserPermissionsRequest) returns (UserPermissionsResponse);
    rpc CheckCourseEnrollment(CourseEnrollmentRequest) returns (CourseEnrollmentResponse);
}

message VideoAccessRequest {
    string user_id = 1;
    string video_id = 2;
    string access_type = 3; // "view", "download", "share"
}

message VideoAccessResponse {
    bool allowed = 1;
    string reason = 2;
    repeated string permissions = 3;
    int64 expires_at = 4;
}
```

### 2. HTTP Integration with Course Service

```go
type CourseServiceClient interface {
    GetCourseByID(courseID uuid.UUID) (*Course, error)
    GetLectureByID(lectureID uuid.UUID) (*Lecture, error)
    UpdateLectureVideo(lectureID, videoID uuid.UUID) error
    GetCourseEnrollments(courseID uuid.UUID) ([]*Enrollment, error)
}
```

### 3. Circuit Breaker Implementation

```go
type VideoServiceHandler struct {
    authCircuitBreaker    *CircuitBreaker
    courseCircuitBreaker  *CircuitBreaker
    cloudflareCircuitBreaker *CircuitBreaker
    redisCircuitBreaker   *CircuitBreaker
}

func (h *VideoServiceHandler) GetVideoWithPermissions(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) (*VideoWithPermissions, error) {
    // Use circuit breaker for auth service call
    authResult, err := h.authCircuitBreaker.Execute(func() (interface{}, error) {
        return h.authService.ValidateVideoAccess(ctx, userID, videoID, "view")
    })

    if err != nil {
        // Fallback to basic permission check
        return h.getVideoWithBasicPermissions(videoID, userID)
    }

    // Continue with full permission validation
    return h.getVideoWithFullPermissions(videoID, userID, authResult)
}
```

This comprehensive technical implementation provides the foundation for a production-grade video streaming platform with advanced network intelligence, real-time adaptation, and detailed analytics capabilities.
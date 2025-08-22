package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB represents a JSONB data type for PostgreSQL
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// Video represents the video metadata table
type Video struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CloudflareUID    string     `json:"cloudflare_uid" gorm:"uniqueIndex;not null"`
	Title            string     `json:"title" gorm:"not null"`
	Description      string     `json:"description"`
	DurationSeconds  *int       `json:"duration_seconds"`
	FileSizeBytes    *int64     `json:"file_size_bytes"`
	UploadUserID     uuid.UUID  `json:"upload_user_id" gorm:"not null"`
	CourseID         *uuid.UUID `json:"course_id"`
	LectureID        *uuid.UUID `json:"lecture_id"`
	Status           string     `json:"status" gorm:"default:'processing'"`
	Visibility       string     `json:"visibility" gorm:"default:'private'"`
	ThumbnailURL     string     `json:"thumbnail_url"`
	StreamURL        string     `json:"stream_url"`
	PreviewURL       string     `json:"preview_url"`
	Metadata         JSONB      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt        time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// VideoQuality represents the video quality variants table
type VideoQuality struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VideoID       uuid.UUID `json:"video_id" gorm:"not null"`
	QualityLabel  string    `json:"quality_label" gorm:"not null"` // '360p', '720p', '1080p'
	BitrateKbps   int       `json:"bitrate_kbps" gorm:"not null"`
	Width         int       `json:"width" gorm:"not null"`
	Height        int       `json:"height" gorm:"not null"`
	FPS           int       `json:"fps" gorm:"default:30"`
	Codec         string    `json:"codec" gorm:"default:'h264'"`
	URL           string    `json:"url" gorm:"not null"`
	FileSizeBytes *int64    `json:"file_size_bytes"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:now()"`
	
	// Relationship
	Video Video `json:"video,omitempty" gorm:"foreignKey:VideoID;references:ID"`
}

// ViewingSession represents the viewing_sessions table
type ViewingSession struct {
	ID                      uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SessionID               string     `json:"session_id" gorm:"not null"`
	UserID                  uuid.UUID  `json:"user_id" gorm:"not null"`
	VideoID                 uuid.UUID  `json:"video_id" gorm:"not null"`
	StartedAt               time.Time  `json:"started_at" gorm:"default:now()"`
	LastHeartbeat           time.Time  `json:"last_heartbeat" gorm:"default:now()"`
	CurrentTimeSeconds      int        `json:"current_time_seconds" gorm:"default:0"`
	CurrentQuality          string     `json:"current_quality"`
	TotalWatchTimeSeconds   int        `json:"total_watch_time_seconds" gorm:"default:0"`
	Completed               bool       `json:"completed" gorm:"default:false"`
	UserAgent               string     `json:"user_agent"`
	IPAddress               string     `json:"ip_address"`
	CreatedAt               time.Time  `json:"created_at" gorm:"default:now()"`
	
	// Relationship
	Video Video `json:"video,omitempty" gorm:"foreignKey:VideoID;references:ID"`
}

// NetworkMetrics represents the network_metrics table
type NetworkMetrics struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SessionID           string    `json:"session_id" gorm:"not null"`
	UserID              uuid.UUID `json:"user_id" gorm:"not null"`
	Timestamp           time.Time `json:"timestamp" gorm:"default:now()"`
	BandwidthMbps       float64   `json:"bandwidth_mbps"`
	LatencyMs           int       `json:"latency_ms"`
	PacketLossPercent   float64   `json:"packet_loss_percent"`
	ConnectionType      string    `json:"connection_type"` // 'wifi', '4g', '5g', 'ethernet'
	QualityScore        int       `json:"quality_score"`   // 1-10 scale
	RecommendedQuality  string    `json:"recommended_quality"`
	BufferHealthSeconds int       `json:"buffer_health_seconds"`
	CreatedAt           time.Time `json:"created_at" gorm:"default:now()"`
}

// VideoPermissions represents the video_permissions table
type VideoPermissions struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VideoID        uuid.UUID  `json:"video_id" gorm:"not null"`
	UserID         *uuid.UUID `json:"user_id"`
	RoleID         *uuid.UUID `json:"role_id"`
	PermissionType string     `json:"permission_type" gorm:"not null"` // 'view', 'download', 'share'
	GrantedBy      uuid.UUID  `json:"granted_by" gorm:"not null"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:now()"`
	
	// Relationship
	Video Video `json:"video,omitempty" gorm:"foreignKey:VideoID;references:ID"`
}

// VideoAnalytics represents the video_analytics table
type VideoAnalytics struct {
	ID                        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VideoID                   uuid.UUID `json:"video_id" gorm:"not null"`
	Date                      time.Time `json:"date" gorm:"not null"`
	TotalViews                int       `json:"total_views" gorm:"default:0"`
	UniqueViewers             int       `json:"unique_viewers" gorm:"default:0"`
	TotalWatchTimeSeconds     int64     `json:"total_watch_time_seconds" gorm:"default:0"`
	AvgWatchTimeSeconds       int       `json:"avg_watch_time_seconds" gorm:"default:0"`
	CompletionRate            float64   `json:"completion_rate" gorm:"default:0"`
	QualityDistribution       JSONB     `json:"quality_distribution" gorm:"type:jsonb"`
	GeographicDistribution    JSONB     `json:"geographic_distribution" gorm:"type:jsonb"`
	DeviceDistribution        JSONB     `json:"device_distribution" gorm:"type:jsonb"`
	CreatedAt                 time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"default:now()"`
	
	// Relationship
	Video Video `json:"video,omitempty" gorm:"foreignKey:VideoID;references:ID"`
}

// Video upload request/response types
type UploadVideoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	CourseID    *uuid.UUID `json:"course_id"`
	LectureID   *uuid.UUID `json:"lecture_id"`
	Visibility  string     `json:"visibility"`
}

type UploadVideoResponse struct {
	VideoID                  uuid.UUID `json:"video_id"`
	CloudflareUID            string    `json:"cloudflare_uid"`
	Title                    string    `json:"title"`
	Status                   string    `json:"status"`
	UploadURL               string    `json:"upload_url"`
	EstimatedProcessingTime string    `json:"estimated_processing_time"`
}

// Session request/response types
type CreateSessionRequest struct {
	DeviceInfo struct {
		UserAgent        string `json:"user_agent"`
		ScreenResolution string `json:"screen_resolution"`
		ConnectionType   string `json:"connection_type"`
	} `json:"device_info"`
}

type CreateSessionResponse struct {
	SessionID          string          `json:"session_id"`
	StreamURL          string          `json:"stream_url"`
	ThumbnailURL       string          `json:"thumbnail_url"`
	Qualities          []VideoQuality  `json:"qualities"`
	RecommendedQuality string          `json:"recommended_quality"`
	WebSocketURL       string          `json:"websocket_url"`
	ExpiresAt          time.Time       `json:"expires_at"`
}

// Network status types
type NetworkStatusUpdate struct {
	BandwidthMbps   float64 `json:"bandwidth_mbps"`
	LatencyMs       int     `json:"latency_ms"`
	PacketLoss      float64 `json:"packet_loss"`
	ConnectionType  string  `json:"connection_type"`
	BufferHealth    int     `json:"buffer_health"`
	CurrentTime     int     `json:"current_time"`
	CurrentQuality  string  `json:"current_quality"`
}

type NetworkStatusResponse struct {
	RecommendedQuality string `json:"recommended_quality"`
	QualityScore       int    `json:"quality_score"`
	ShouldPreload      bool   `json:"should_preload"`
	BufferTarget       int    `json:"buffer_target"`
}

// WebSocket message types
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSNetworkStatus struct {
	SessionID       string  `json:"session_id"`
	UserID          string  `json:"user_id"`
	Timestamp       string  `json:"timestamp"`
	BandwidthMbps   float64 `json:"bandwidth_mbps"`
	LatencyMs       int     `json:"latency_ms"`
	PacketLoss      float64 `json:"packet_loss"`
	ConnectionType  string  `json:"connection_type"`
	QualityScore    int     `json:"quality_score"`
	CurrentQuality  string  `json:"current_quality"`
	RecommendedQuality string `json:"recommended_quality"`
	BufferHealth    int     `json:"buffer_health"`
}

type WSQualityChange struct {
	SessionID   string `json:"session_id"`
	VideoID     string `json:"video_id"`
	FromQuality string `json:"from_quality"`
	ToQuality   string `json:"to_quality"`
	Reason      string `json:"reason"`
	Timestamp   string `json:"timestamp"`
}

type WSQualityRecommendation struct {
	RecommendedQuality string  `json:"recommended_quality"`
	Reason            string  `json:"reason"`
	Confidence        float64 `json:"confidence"`
}

type WSPreloadInstruction struct {
	Segments []string `json:"segments"`
	Priority string   `json:"priority"`
}

type WSAnalyticsEvent struct {
	Event string `json:"event"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}
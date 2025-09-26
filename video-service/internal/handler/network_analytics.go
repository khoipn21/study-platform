package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"video-service/internal/model"
	"video-service/internal/service"
)

type NetworkAnalyticsHandler struct {
	networkService *service.NetworkIntelligenceService
	videoService   *service.VideoService
}

func NewNetworkAnalyticsHandler(networkService *service.NetworkIntelligenceService, videoService *service.VideoService) *NetworkAnalyticsHandler {
	return &NetworkAnalyticsHandler{
		networkService: networkService,
		videoService:   videoService,
	}
}

// GetNetworkStatus handles network status retrieval for a session
// GET /api/videos/sessions/{session_id}/network-status
func (nah *NetworkAnalyticsHandler) GetNetworkStatus(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Get cached network metrics
	metrics, err := nah.networkService.GetCachedNetworkMetrics(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No network data available for this session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":           sessionID,
		"bandwidth_mbps":       metrics.BandwidthMbps,
		"latency_ms":           metrics.LatencyMs,
		"packet_loss_percent":  metrics.PacketLossPercent,
		"connection_type":      metrics.ConnectionType,
		"quality_score":        metrics.QualityScore,
		"recommended_quality":  metrics.RecommendedQuality,
		"buffer_health_seconds": metrics.BufferHealthSeconds,
		"timestamp":            metrics.Timestamp,
	})
}

// GetQualityRecommendation handles quality recommendation requests
// GET /api/videos/sessions/{session_id}/quality-recommendation
func (nah *NetworkAnalyticsHandler) GetQualityRecommendation(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Get connection type from query parameter
	connectionType := c.DefaultQuery("connection_type", "wifi")

	// Get quality recommendation
	recommendedQuality := nah.networkService.GetQualityRecommendation(c.Request.Context(), connectionType)

	// Try to get cached recommendation if available
	if cached, err := nah.networkService.GetCachedQualityRecommendation(c.Request.Context(), sessionID); err == nil && cached != "" {
		recommendedQuality = cached
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":           sessionID,
		"recommended_quality":  recommendedQuality,
		"connection_type":      connectionType,
		"timestamp":            time.Now().Format(time.RFC3339),
	})
}

// GetNetworkPattern analyzes network patterns for a session
// GET /api/videos/sessions/{session_id}/network-pattern
func (nah *NetworkAnalyticsHandler) GetNetworkPattern(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Parse window parameter (default: 5 minutes)
	windowStr := c.DefaultQuery("window_minutes", "5")
	windowMinutes, err := strconv.Atoi(windowStr)
	if err != nil || windowMinutes <= 0 {
		windowMinutes = 5
	}

	pattern, err := nah.networkService.AnalyzeNetworkPattern(c.Request.Context(), sessionID, windowMinutes)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pattern)
}

// PostNetworkDiagnostics performs network diagnostics
// POST /api/videos/sessions/{session_id}/network-diagnostics
func (nah *NetworkAnalyticsHandler) PostNetworkDiagnostics(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var req struct {
		TestDurationSeconds int    `json:"test_duration_seconds"`
		TestType           string `json:"test_type"` // "bandwidth", "latency", "packet_loss", "full"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.TestDurationSeconds == 0 {
		req.TestDurationSeconds = 10
	}
	if req.TestType == "" {
		req.TestType = "full"
	}

	// For now, simulate diagnostic results
	// In a real implementation, this would trigger actual network tests
	diagnostics := gin.H{
		"session_id":     sessionID,
		"test_type":      req.TestType,
		"test_duration":  req.TestDurationSeconds,
		"status":         "completed",
		"results": gin.H{
			"bandwidth_test": gin.H{
				"download_mbps": 25.5,
				"upload_mbps":   12.3,
				"stability":     "good",
			},
			"latency_test": gin.H{
				"average_ms": 45,
				"min_ms":     32,
				"max_ms":     78,
				"jitter_ms":  12,
			},
			"packet_loss_test": gin.H{
				"loss_percent":    0.02,
				"packets_sent":    1000,
				"packets_lost":    0,
				"test_reliable":   true,
			},
			"quality_assessment": gin.H{
				"overall_score":        8,
				"recommended_quality":  "720p",
				"buffer_recommendation": 10,
				"preload_enabled":      true,
			},
		},
		"timestamp":      time.Now().Format(time.RFC3339),
		"recommendations": []string{
			"Network conditions are good for 720p streaming",
			"Preloading enabled for optimal experience",
			"Consider 1080p if network remains stable",
		},
	}

	c.JSON(http.StatusOK, diagnostics)
}

// GetSessionAnalytics provides comprehensive session analytics
// GET /api/videos/sessions/{session_id}/analytics
func (nah *NetworkAnalyticsHandler) GetSessionAnalytics(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// In a real implementation, this would query historical data
	// For now, provide comprehensive mock analytics
	analytics := gin.H{
		"session_id": sessionID,
		"overview": gin.H{
			"session_duration_seconds":    3600,
			"total_watch_time_seconds":   2850,
			"completion_rate":            0.79,
			"quality_changes":            12,
			"buffer_events":              3,
			"network_interruptions":      1,
		},
		"network_performance": gin.H{
			"average_bandwidth_mbps":     8.5,
			"average_latency_ms":         65,
			"average_packet_loss":        0.015,
			"stability_score":            7,
			"connection_type_history": gin.H{
				"wifi":     85,
				"4g":       15,
				"ethernet": 0,
			},
		},
		"quality_distribution": gin.H{
			"1080p": gin.H{"percentage": 45, "duration_seconds": 1282},
			"720p":  gin.H{"percentage": 35, "duration_seconds": 997},
			"480p":  gin.H{"percentage": 15, "duration_seconds": 427},
			"360p":  gin.H{"percentage": 5, "duration_seconds": 142},
		},
		"performance_timeline": []gin.H{
			{
				"timestamp":           time.Now().Add(-1*time.Hour).Format(time.RFC3339),
				"quality":             "1080p",
				"bandwidth_mbps":      12.5,
				"latency_ms":          45,
				"buffer_health":       15,
				"quality_score":       9,
			},
			{
				"timestamp":           time.Now().Add(-45*time.Minute).Format(time.RFC3339),
				"quality":             "720p",
				"bandwidth_mbps":      6.2,
				"latency_ms":          85,
				"buffer_health":       8,
				"quality_score":       6,
			},
		},
		"recommendations": gin.H{
			"overall_experience":    "good",
			"suggested_improvements": []string{
				"Consider upgrading connection for consistent 1080p",
				"Buffer size optimization implemented successfully",
				"Network stability is acceptable for streaming",
			},
			"optimal_settings": gin.H{
				"preferred_quality":   "720p",
				"buffer_target":       12,
				"preload_segments":    3,
			},
		},
		"generated_at": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, analytics)
}

// UpdateBandwidthEstimate handles bandwidth estimation updates
// POST /api/videos/sessions/{session_id}/bandwidth-estimate
func (nah *NetworkAnalyticsHandler) UpdateBandwidthEstimate(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var req struct {
		EstimatedBandwidthMbps float64 `json:"estimated_bandwidth_mbps" binding:"required"`
		MeasurementMethod      string  `json:"measurement_method"`
		Confidence            float64 `json:"confidence"`
		TestDurationSeconds   int     `json:"test_duration_seconds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create network update from bandwidth estimate
	networkUpdate := &model.NetworkStatusUpdate{
		BandwidthMbps:   req.EstimatedBandwidthMbps,
		ConnectionType:  "unknown", // Will be determined by client
		CurrentTime:     0,         // Not relevant for bandwidth update
	}

	// Process through network intelligence service
	response, err := nah.networkService.ProcessNetworkUpdate(c.Request.Context(), sessionID, networkUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{
		"session_id":                sessionID,
		"received_bandwidth_mbps":   req.EstimatedBandwidthMbps,
		"measurement_method":        req.MeasurementMethod,
		"confidence":                req.Confidence,
		"recommended_quality":       response.RecommendedQuality,
		"quality_score":             response.QualityScore,
		"should_preload":           response.ShouldPreload,
		"buffer_target":            response.BufferTarget,
		"timestamp":                time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, result)
}

// GetActiveNetworkSessions returns all active network monitoring sessions
// GET /api/videos/network/active-sessions
func (nah *NetworkAnalyticsHandler) GetActiveNetworkSessions(c *gin.Context) {
	// Get user ID from context if available for filtering
	var userID *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsedUserID, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsedUserID
		}
	}

	// In a real implementation, this would query active sessions from Redis
	// For now, provide mock data
	activeSessions := []gin.H{
		{
			"session_id":              "sess_12345",
			"user_id":                 userID,
			"video_id":                "video_67890",
			"current_quality":         "720p",
			"bandwidth_mbps":          8.5,
			"latency_ms":              65,
			"connection_type":         "wifi",
			"quality_score":           7,
			"buffer_health_seconds":   12,
			"last_update":             time.Now().Add(-30*time.Second).Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"active_sessions": activeSessions,
		"total_count":     len(activeSessions),
		"timestamp":       time.Now().Format(time.RFC3339),
	})
}

// GetNetworkMetricsHistory retrieves historical network metrics
// GET /api/videos/sessions/{session_id}/network-history
func (nah *NetworkAnalyticsHandler) GetNetworkMetricsHistory(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Parse time range parameters
	hoursStr := c.DefaultQuery("hours", "1")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		hours = 1
	}

	// In a real implementation, this would query the database
	// For now, generate mock historical data
	history := []gin.H{}
	now := time.Now()

	// Generate data points every 5 minutes for the specified hours
	intervalMinutes := 5
	totalPoints := (hours * 60) / intervalMinutes

	for i := 0; i < totalPoints; i++ {
		timestamp := now.Add(time.Duration(-i*intervalMinutes) * time.Minute)

		// Simulate varying network conditions
		baseQuality := 7
		if i > totalPoints/2 {
			baseQuality = 5 // Simulate network degradation
		}

		history = append(history, gin.H{
			"timestamp":               timestamp.Format(time.RFC3339),
			"bandwidth_mbps":          float64(baseQuality) + 1.5,
			"latency_ms":              100 - baseQuality*10,
			"packet_loss_percent":     float64(10-baseQuality) / 1000,
			"connection_type":         "wifi",
			"quality_score":           baseQuality,
			"recommended_quality":     map[int]string{5: "480p", 6: "720p", 7: "720p", 8: "1080p"}[baseQuality],
			"buffer_health_seconds":   baseQuality + 5,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":      sessionID,
		"hours_requested": hours,
		"data_points":     len(history),
		"history":         history,
		"generated_at":    now.Format(time.RFC3339),
	})
}
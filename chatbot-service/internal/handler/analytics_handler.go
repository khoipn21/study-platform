package handler

import (
	"net/http"
	"strconv"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsServiceInterface
}

func NewAnalyticsHandler(analyticsService service.AnalyticsServiceInterface) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetOverallAnalytics handles GET /analytics/overall
func (h *AnalyticsHandler) GetOverallAnalytics(c *gin.Context) {
	req := h.parseAnalyticsRequest(c)
	
	analytics, err := h.analyticsService.GetOverallAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetUserAnalytics handles GET /analytics/user/:userID
func (h *AnalyticsHandler) GetUserAnalytics(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	req := h.parseAnalyticsRequest(c)
	
	analytics, err := h.analyticsService.GetUserAnalytics(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetMyAnalytics handles GET /analytics/me (authenticated user's analytics)
func (h *AnalyticsHandler) GetMyAnalytics(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	req := h.parseAnalyticsRequest(c)
	
	analytics, err := h.analyticsService.GetUserAnalytics(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetCourseAnalytics handles GET /analytics/course/:courseID
func (h *AnalyticsHandler) GetCourseAnalytics(c *gin.Context) {
	courseIDParam := c.Param("courseID")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	req := h.parseAnalyticsRequest(c)
	
	analytics, err := h.analyticsService.GetCourseAnalytics(c.Request.Context(), courseID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetTimeBasedAnalytics handles GET /analytics/time-based
func (h *AnalyticsHandler) GetTimeBasedAnalytics(c *gin.Context) {
	req := h.parseAnalyticsRequest(c)
	
	analytics, err := h.analyticsService.GetTimeBasedAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve time-based analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetRealTimeMetrics handles GET /analytics/real-time
func (h *AnalyticsHandler) GetRealTimeMetrics(c *gin.Context) {
	metrics, err := h.analyticsService.GetRealTimeMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real-time metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetSessionMetrics handles GET /analytics/sessions
func (h *AnalyticsHandler) GetSessionMetrics(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	sessions, err := h.analyticsService.GetSessionMetrics(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// GetUsageStats handles GET /analytics/usage
func (h *AnalyticsHandler) GetUsageStats(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	period := c.DefaultQuery("period", "monthly")
	if period != "daily" && period != "weekly" && period != "monthly" {
		period = "monthly"
	}

	stats, err := h.analyticsService.GetUsageStats(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve usage stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GenerateReport handles POST /analytics/report
func (h *AnalyticsHandler) GenerateReport(c *gin.Context) {
	var req model.AnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails, try to parse from query parameters
		req = *h.parseAnalyticsRequest(c)
	}

	report, err := h.analyticsService.GenerateAnalyticsReport(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate analytics report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetResponseQuality handles GET /analytics/quality
func (h *AnalyticsHandler) GetResponseQuality(c *gin.Context) {
	req := h.parseAnalyticsRequest(c)
	
	quality, err := h.analyticsService.CalculateResponseQuality(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate response quality"})
		return
	}

	c.JSON(http.StatusOK, quality)
}

// Helper method to parse analytics request from query parameters
func (h *AnalyticsHandler) parseAnalyticsRequest(c *gin.Context) *model.AnalyticsRequest {
	req := &model.AnalyticsRequest{}

	// Parse dates
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = &endDate
		}
	}

	// Parse UUIDs
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			req.UserID = &userID
		}
	}

	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			req.CourseID = &courseID
		}
	}

	// Parse period and granularity
	req.Period = c.DefaultQuery("period", "daily")
	req.Granularity = c.DefaultQuery("granularity", "day")

	// Validate period
	if req.Period != "daily" && req.Period != "weekly" && req.Period != "monthly" {
		req.Period = "daily"
	}

	return req
}

// GetAnalyticsDashboard handles GET /analytics/dashboard
func (h *AnalyticsHandler) GetAnalyticsDashboard(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get real-time metrics
	realTimeMetrics, err := h.analyticsService.GetRealTimeMetrics(c.Request.Context())
	if err != nil {
		realTimeMetrics = &model.RealTimeMetrics{} // fallback to empty metrics
	}

	// Get user analytics for the last 30 days
	req := &model.AnalyticsRequest{
		Period: "daily",
		UserID: &userID,
	}
	userAnalytics, err := h.analyticsService.GetUserAnalytics(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dashboard data"})
		return
	}

	// Get usage stats
	usageStats, err := h.analyticsService.GetUsageStats(c.Request.Context(), userID, "monthly")
	if err != nil {
		usageStats = &model.UsageStats{} // fallback to empty stats
	}

	// Get recent session metrics
	sessionMetrics, err := h.analyticsService.GetSessionMetrics(c.Request.Context(), userID, 5)
	if err != nil {
		sessionMetrics = []*model.SessionMetrics{} // fallback to empty array
	}

	dashboard := gin.H{
		"real_time_metrics": realTimeMetrics,
		"user_analytics":    userAnalytics.Data,
		"usage_stats":       usageStats,
		"recent_sessions":   sessionMetrics,
		"generated_at":      time.Now(),
	}

	c.JSON(http.StatusOK, dashboard)
}
package handler

import (
	"net/http"
	"time"

	"instructor-dashboard-service/internal/model"
	"instructor-dashboard-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetRevenueAnalytics handles GET /api/v1/instructor/analytics/revenue
func (h *AnalyticsHandler) GetRevenueAnalytics(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	analytics, err := h.analyticsService.GetRevenueAnalytics(instructorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get revenue analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetEngagementAnalytics handles GET /api/v1/instructor/analytics/engagement
func (h *AnalyticsHandler) GetEngagementAnalytics(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	analytics, err := h.analyticsService.GetEngagementAnalytics(instructorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get engagement analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetStudentAnalytics handles GET /api/v1/instructor/analytics/students
func (h *AnalyticsHandler) GetStudentAnalytics(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	analytics, err := h.analyticsService.GetStudentAnalytics(instructorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetCourseAnalytics handles GET /api/v1/instructor/courses/:id/analytics
func (h *AnalyticsHandler) GetCourseAnalytics(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	analytics, err := h.analyticsService.GetCourseAnalytics(instructorID, courseID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetVideoAnalytics handles GET /api/v1/instructor/videos/analytics
func (h *AnalyticsHandler) GetVideoAnalytics(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	analytics, err := h.analyticsService.GetVideoAnalytics(instructorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetVideoEngagement handles GET /api/v1/instructor/videos/:id/engagement
func (h *AnalyticsHandler) GetVideoEngagement(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	req := parseAnalyticsRequest(c)

	engagement, err := h.analyticsService.GetVideoEngagement(instructorID, videoID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video engagement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    engagement,
	})
}

// GetAnalyticsSummary handles GET /api/v1/instructor/analytics/summary
func (h *AnalyticsHandler) GetAnalyticsSummary(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "monthly"
	}

	summary, err := h.analyticsService.GetAnalyticsSummary(instructorID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// parseAnalyticsRequest parses common analytics request parameters
func parseAnalyticsRequest(c *gin.Context) *model.AnalyticsRequest {
	req := &model.AnalyticsRequest{
		Period:   c.Query("period"),
		Timezone: c.Query("timezone"),
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	if req.Timezone == "" {
		req.Timezone = "UTC"
	}

	// Parse start date (support both camelCase and snake_case)
	startDateStr := c.Query("start_date")
	if startDateStr == "" {
		startDateStr = c.Query("startDate")
	}
	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = &startDate
		}
	}

	// Parse end date (support both camelCase and snake_case)
	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		endDateStr = c.Query("endDate")
	}
	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = &endDate
		}
	}

	// Parse granularity (for revenue analytics)
	granularity := c.Query("granularity")
	if granularity != "" {
		req.Period = granularity
	}

	// Parse course ID if provided
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			req.CourseID = &courseID
		}
	}

	// Parse video ID if provided
	if videoIDStr := c.Query("video_id"); videoIDStr != "" {
		if videoID, err := uuid.Parse(videoIDStr); err == nil {
			req.VideoID = &videoID
		}
	}

	return req
}
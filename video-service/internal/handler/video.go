package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"video-service/internal/model"
	"video-service/internal/service"
)

type VideoHandler struct {
	videoService *service.VideoService
}

func NewVideoHandler(videoService *service.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

// GetUploadURL handles getting upload URL for video upload
// POST /api/videos/upload-url
func (vh *VideoHandler) GetUploadURL(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Filename string `json:"filename" binding:"required"`
		Size     int64  `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := vh.videoService.GetUploadURL(c.Request.Context(), userID, req.Filename, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UploadVideo handles video upload requests
// POST /api/videos/upload
func (vh *VideoHandler) UploadVideo(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req model.UploadVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := vh.videoService.UploadVideo(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetVideo handles get video requests
// GET /api/videos/{video_id}
func (vh *VideoHandler) GetVideo(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get user ID from context if available
	var userID *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsedUserID, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsedUserID
		}
	}

	video, err := vh.videoService.GetVideo(c.Request.Context(), videoID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

// CreateSession handles create viewing session requests
// POST /api/videos/{video_id}/sessions
func (vh *VideoHandler) CreateSession(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req model.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract client IP address for session tracking
	clientIP := vh.getClientIP(c)

	response, err := vh.videoService.CreateViewingSession(c.Request.Context(), videoID, userID, &req, clientIP)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if err.Error() == "video is not ready for streaming" {
			c.JSON(http.StatusConflict, gin.H{"error": "Video is not ready for streaming"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateNetworkStatus handles network status updates
// POST /api/videos/sessions/{session_id}/network
func (vh *VideoHandler) UpdateNetworkStatus(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var update model.NetworkStatusUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := vh.videoService.UpdateNetworkStatus(c.Request.Context(), sessionID, &update)
	if err != nil {
		if err.Error() == "session not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSessionProgress handles session progress updates
// PUT /api/videos/sessions/{session_id}/progress
func (vh *VideoHandler) UpdateSessionProgress(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var req struct {
		CurrentTime int    `json:"current_time" binding:"required"`
		Quality     string `json:"quality"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vh.videoService.UpdateSessionProgress(c.Request.Context(), sessionID, req.CurrentTime, req.Quality)
	if err != nil {
		if err.Error() == "viewing session not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}

// ListUserVideos handles listing user's videos
// GET /api/videos/user/{user_id}
func (vh *VideoHandler) ListUserVideos(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	videos, err := vh.videoService.ListUserVideos(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"limit":  limit,
		"offset": offset,
	})
}

// SearchVideos handles video search requests
// GET /api/videos/search
func (vh *VideoHandler) SearchVideos(c *gin.Context) {
	query := c.Query("q")
	// Allow empty queries - repository will return all public videos

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	videos, err := vh.videoService.SearchVideos(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"query":  query,
		"limit":  limit,
		"offset": offset,
	})
}

// DeleteVideo handles video deletion requests
// DELETE /api/videos/{video_id}
func (vh *VideoHandler) DeleteVideo(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = vh.videoService.DeleteVideo(c.Request.Context(), videoID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if err.Error() == "video not found or already deleted" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

// CloudflareWebhook handles webhooks from Cloudflare Stream
// POST /api/videos/webhooks/cloudflare
func (vh *VideoHandler) CloudflareWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: Add webhook signature verification here
	// For now, we'll accept all webhooks for development

	// Debug logging
	if uid, ok := webhookData["uid"].(string); ok {
		c.Header("X-Debug-UID", uid)
	}

	err := vh.videoService.ProcessCloudflareWebhook(c.Request.Context(), webhookData)
	if err != nil {
		// Add debug info to error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"debug": "Error in ProcessCloudflareWebhook",
			"uid":   webhookData["uid"],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

// UpdateVideoStatus manually updates video status (for recovery)
// PUT /api/videos/{video_id}/status
func (vh *VideoHandler) UpdateVideoStatus(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"uploading":  true,
		"processing": true,
		"ready":      true,
		"error":      true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid statuses: uploading, processing, ready, error"})
		return
	}

	err = vh.videoService.UpdateVideoStatus(c.Request.Context(), videoID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video status updated successfully"})
}

// GetVideoAnalytics handles video analytics requests
// GET /api/videos/{video_id}/analytics
func (vh *VideoHandler) GetVideoAnalytics(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	period := c.DefaultQuery("period", "7d")

	// TODO: Implement analytics retrieval
	// For now, return mock data
	analytics := gin.H{
		"video_id":       videoID,
		"period":         period,
		"total_views":    1250,
		"unique_viewers": 890,
		"total_watch_time": 156000,
		"avg_watch_time": 175,
		"completion_rate": 0.65,
		"daily_stats": []gin.H{
			{
				"date":            "2024-01-01",
				"views":           180,
				"unique_viewers":  150,
				"watch_time":      22500,
			},
		},
		"quality_distribution": gin.H{
			"360p":  25,
			"720p":  55,
			"1080p": 20,
		},
	}

	c.JSON(http.StatusOK, analytics)
}

// ListCourseVideos handles listing videos for a course
// GET /api/videos/course/{course_id}
func (vh *VideoHandler) ListCourseVideos(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	videos, err := vh.videoService.ListCourseVideos(c.Request.Context(), courseID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos":    videos,
		"course_id": courseID,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateVideo handles video update requests
// PUT /api/videos/{video_id}
func (vh *VideoHandler) UpdateVideo(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateReq struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get video first to check permissions
	video, err := vh.videoService.GetVideo(c.Request.Context(), videoID, &userID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user owns the video
	if video.UploadUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update video fields
	if updateReq.Title != "" {
		video.Title = updateReq.Title
	}
	if updateReq.Description != "" {
		video.Description = updateReq.Description
	}
	if updateReq.Visibility != "" {
		video.Visibility = updateReq.Visibility
	}

	// TODO: Implement video update in service
	// For now, return success
	c.JSON(http.StatusOK, video)
}

// getClientIP extracts the real client IP address from various headers
func (vh *VideoHandler) getClientIP(c *gin.Context) string {
	// Check for X-Forwarded-For header (most common for proxies)
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := strings.Index(forwarded, ","); idx != -1 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}

	// Check for X-Real-IP header
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// Check for CF-Connecting-IP (Cloudflare)
	if cfIP := c.GetHeader("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// Fallback to RemoteAddr
	ip := c.ClientIP()
	if ip == "" {
		ip = "127.0.0.1" // Default fallback for empty IP
	}
	return ip
}
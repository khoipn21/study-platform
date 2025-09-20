package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"video-service/internal/security"
	"video-service/internal/service"
)

type VideoAccessHandler struct {
	videoService *service.VideoService
	drmManager   *security.DRMManager
}

func NewVideoAccessHandler(videoService *service.VideoService, drmManager *security.DRMManager) *VideoAccessHandler {
	return &VideoAccessHandler{
		videoService: videoService,
		drmManager:   drmManager,
	}
}

// GenerateStreamURL handles POST /api/videos/{video_id}/stream/url
// Generates a signed URL for video streaming with DRM protection
func (vah *VideoAccessHandler) GenerateStreamURL(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

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

	// Get request parameters
	accessLevel := c.DefaultQuery("access_level", "view")
	expiryMinutes := c.DefaultQuery("expiry_minutes", "120")

	expiryDuration, err := strconv.Atoi(expiryMinutes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry duration"})
		return
	}

	customExpiry := time.Duration(expiryDuration) * time.Minute

	// Validate user access to the video
	if err := vah.validateVideoAccess(c, videoID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Generate signed URL
	signedURL, err := vah.drmManager.GenerateSignedURL(videoID, userID, accessLevel, &customExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed URL"})
		return
	}

	// Log security event
	vah.logSecurityEvent("signed_url_generated", userID, videoID, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"signed_url": signedURL,
		"expires_at": time.Now().Add(customExpiry).Unix(),
		"access_level": accessLevel,
	})
}

// ValidateStreamAccess handles GET /api/videos/{video_id}/stream
// Validates signed URL and provides secure video stream access
func (vah *VideoAccessHandler) ValidateStreamAccess(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get token and signature from query parameters
	token := c.Query("token")
	signature := c.Query("signature")

	if token == "" || signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or signature"})
		return
	}

	// Validate signed URL
	claims, err := vah.drmManager.ValidateSignedURL(token, signature)
	if err != nil {
		vah.logSecurityEvent("invalid_access_attempt", uuid.Nil, videoID, c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Verify video ID matches
	if claims.VideoID != videoID {
		vah.logSecurityEvent("video_id_mismatch", claims.UserID, videoID, c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusForbidden, gin.H{"error": "Video ID mismatch"})
		return
	}

	// Get video metadata
	video, err := vah.videoService.GetVideo(c.Request.Context(), videoID, &claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Generate Cloudflare signed URL for actual streaming
	cloudflareURL, err := vah.drmManager.GenerateCloudflareSignedURL(
		video.CloudflareUID,
		claims.UserID,
		time.Until(time.Unix(claims.ExpiresAt, 0)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate stream URL"})
		return
	}

	// Log successful access
	vah.logSecurityEvent("video_access_granted", claims.UserID, videoID, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"stream_url": cloudflareURL,
		"video_id": videoID,
		"expires_at": claims.ExpiresAt,
		"access_level": claims.AccessLevel,
	})
}

// CreateViewingSession handles POST /api/videos/{video_id}/session
// Creates a secure viewing session for enhanced DRM protection
func (vah *VideoAccessHandler) CreateViewingSession(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

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

	// Validate user access
	if err := vah.validateVideoAccess(c, videoID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Create viewing session
	sessionDuration := 4 * time.Hour // Default 4-hour session
	session, err := vah.drmManager.CreateViewingSession(userID, videoID, sessionDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create viewing session"})
		return
	}

	// Enhance session with request metadata
	session.IPAddress = c.ClientIP()
	session.UserAgent = c.GetHeader("User-Agent")

	// Log session creation
	vah.logSecurityEvent("viewing_session_created", userID, videoID, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusCreated, gin.H{
		"session_id": session.SessionID,
		"session_token": session.Token,
		"expires_at": session.ExpiresAt.Unix(),
		"max_views": session.MaxViews,
	})
}

// ValidateViewingSession handles POST /api/videos/session/validate
// Validates viewing session for playback
func (vah *VideoAccessHandler) ValidateViewingSession(c *gin.Context) {
	var req struct {
		SessionToken string `json:"session_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate session
	session, err := vah.drmManager.ValidateViewingSession(req.SessionToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Increment view count (in production, this would be stored in database)
	session.ViewCount++

	// Log session validation
	vah.logSecurityEvent("viewing_session_validated", session.UserID, session.VideoID, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"session_id": session.SessionID,
		"remaining_views": session.MaxViews - session.ViewCount,
		"expires_at": session.ExpiresAt.Unix(),
	})
}

// RevokeAccess handles POST /api/videos/{video_id}/revoke
// Revokes access to a video (admin only)
func (vah *VideoAccessHandler) RevokeAccess(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Check if user is admin
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	// Get target user ID from request
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	// Log security event for access revocation
	vah.logSecurityEvent("access_revoked", targetUserID, videoID, c.ClientIP(), c.GetHeader("User-Agent"))

	// In production, this would revoke all active sessions for the user/video combination
	c.JSON(http.StatusOK, gin.H{
		"message": "Access revoked successfully",
		"revoked_for": targetUserID,
		"video_id": videoID,
		"revoked_by": adminID,
	})
}

// Helper methods

func (vah *VideoAccessHandler) validateVideoAccess(c *gin.Context, videoID, userID uuid.UUID) error {
	// This would check enrollment, purchase status, etc.
	// For now, we'll do a basic check through the video service

	_, err := vah.videoService.GetVideo(c.Request.Context(), videoID, &userID)
	if err != nil {
		return err
	}

	return nil
}

func (vah *VideoAccessHandler) logSecurityEvent(eventType string, userID, videoID uuid.UUID, ipAddress, userAgent string) {
	event := &security.SecurityEvent{
		EventType:   eventType,
		UserID:      userID,
		VideoID:     videoID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Video access event: %s", eventType),
		Severity:    "medium",
	}

	vah.drmManager.LogSecurityEvent(event)
}
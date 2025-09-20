package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"video-service/internal/middleware"
	"video-service/internal/websocket"
)

type WebSocketHandler struct {
	hub        *websocket.Hub
	authMW     *middleware.WebSocketAuthMiddleware
	rateLimiter *middleware.WebSocketRateLimiter
}

func NewWebSocketHandler(hub *websocket.Hub, authMW *middleware.WebSocketAuthMiddleware) *WebSocketHandler {
	// Create rate limiter: max 5 connections per user per 5 minutes
	rateLimiter := middleware.NewWebSocketRateLimiter(5, 5*time.Minute)

	return &WebSocketHandler{
		hub:         hub,
		authMW:      authMW,
		rateLimiter: rateLimiter,
	}
}

// HandleWebSocket handles WebSocket connections with secure authentication
// GET /api/videos/ws/{session_id}
func (wh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session ID is required",
			"code":  "MISSING_SESSION_ID",
		})
		return
	}

	// Parse session ID
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID format",
			"code":  "INVALID_SESSION_ID",
		})
		return
	}

	// Authenticate WebSocket upgrade request
	claims, err := wh.authMW.AuthenticateWebSocketUpgrade(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "WebSocket authentication failed: " + err.Error(),
			"code":  "WS_AUTH_FAILED",
		})

		// Log failed authentication attempt
		wh.logSecurityEvent("ws_auth_failed", uuid.Nil, uuid.Nil, c.ClientIP(), err.Error())
		return
	}

	// Verify session ID matches token
	if claims.SessionID != sessionUUID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Session ID mismatch",
			"code":  "SESSION_ID_MISMATCH",
		})

		wh.logSecurityEvent("ws_session_mismatch", claims.UserID, claims.VideoID, c.ClientIP(), "Session ID in token does not match URL parameter")
		return
	}

	// Check rate limiting
	if !wh.rateLimiter.AllowConnection(claims.UserID.String()) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded for WebSocket connections",
			"code":  "WS_RATE_LIMIT_EXCEEDED",
		})

		wh.logSecurityEvent("ws_rate_limit", claims.UserID, claims.VideoID, c.ClientIP(), "Too many WebSocket connection attempts")
		return
	}

	// Validate video streaming scope
	if err := wh.authMW.ValidateWebSocketAccess(claims, "video_stream", claims.VideoID.String()); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions: " + err.Error(),
			"code":  "WS_INSUFFICIENT_PERMISSIONS",
		})
		return
	}

	// Log successful connection
	wh.logSecurityEvent("ws_connect_success", claims.UserID, claims.VideoID, c.ClientIP(), "WebSocket connection established")

	// Serve secure WebSocket connection
	websocket.ServeWS(wh.hub, c.Writer, c.Request, claims.SessionID.String(), claims.UserID.String(), claims.VideoID.String())
}

// CreateWebSocketToken creates a WebSocket authentication token
// POST /api/videos/ws/token
func (wh *WebSocketHandler) CreateWebSocketToken(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id" binding:"required"`
		SessionID string `json:"session_id" binding:"required"`
		VideoID   string `json:"video_id" binding:"required"`
		UserRole  string `json:"user_role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	videoID, err := uuid.Parse(req.VideoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Create WebSocket token
	token, err := wh.authMW.CreateSessionWebSocketToken(userID, sessionID, videoID, req.UserRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create WebSocket token",
			"code":  "TOKEN_CREATION_FAILED",
		})
		return
	}

	// Log token creation
	wh.logSecurityEvent("ws_token_created", userID, videoID, c.ClientIP(), "WebSocket token created")

	c.JSON(http.StatusOK, gin.H{
		"ws_token": token,
		"expires_in": 4 * 60 * 60, // 4 hours in seconds
		"scopes": wh.getScopesForRole(req.UserRole),
	})
}

// logSecurityEvent logs security events for monitoring
func (wh *WebSocketHandler) logSecurityEvent(eventType string, userID, videoID uuid.UUID, ipAddress, description string) {
	log := &middleware.WebSocketConnectionLog{
		UserID:      userID,
		VideoID:     videoID,
		IPAddress:   ipAddress,
		ConnectedAt: time.Now(),
		EventType:   eventType,
		ErrorReason: description,
	}
	wh.authMW.LogWebSocketConnection(log)
}

// getScopesForRole returns the scopes for a given user role
func (wh *WebSocketHandler) getScopesForRole(role string) []string {
	switch role {
	case "admin":
		return []string{"video_stream", "quality_control", "analytics", "broadcast", "admin"}
	case "instructor":
		return []string{"video_stream", "quality_control", "analytics", "broadcast"}
	default:
		return []string{"video_stream"}
	}
}

// GetWebSocketStats returns WebSocket connection statistics
// GET /api/videos/ws/stats
func (wh *WebSocketHandler) GetWebSocketStats(c *gin.Context) {
	stats := gin.H{
		"active_connections": wh.hub.GetActiveClients(),
		"timestamp":         c.Query("timestamp"),
	}

	c.JSON(http.StatusOK, stats)
}

// BroadcastMessage sends a message to all connected clients
// POST /api/videos/ws/broadcast
func (wh *WebSocketHandler) BroadcastMessage(c *gin.Context) {
	var message struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create WebSocket message
	wsMessage := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: message.Type,
		Data: message.Data,
	}

	// Broadcast to all clients
	wh.hub.BroadcastToAll(wsMessage)

	c.JSON(http.StatusOK, gin.H{"message": "Broadcast sent successfully"})
}

// SendToSession sends a message to a specific session
// POST /api/videos/ws/session/{session_id}/send
func (wh *WebSocketHandler) SendToSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	var message struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create WebSocket message
	wsMessage := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: message.Type,
		Data: message.Data,
	}

	// Send to specific session
	wh.hub.SendToSession(sessionID, wsMessage)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// GetSessionInfo returns information about a specific session
// GET /api/videos/ws/session/{session_id}
func (wh *WebSocketHandler) GetSessionInfo(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	client := wh.hub.GetSessionClient(sessionID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	info := client.GetConnectionInfo()
	c.JSON(http.StatusOK, info)
}
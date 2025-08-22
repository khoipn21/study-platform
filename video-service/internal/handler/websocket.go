package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"video-service/internal/websocket"
)

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket handles WebSocket connections
// GET /api/videos/ws/{session_id}
func (wh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Get user ID from query parameter or JWT token
	userID := c.Query("user_id")
	if userID == "" {
		// Try to get from context if available
		if userIDFromContext, exists := c.Get("user_id"); exists {
			userID = userIDFromContext.(string)
		}
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is required"})
		return
	}

	// Get video ID from query parameter
	videoID := c.Query("video_id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	// Optional JWT token validation can be done here
	token := c.Query("token")
	if token == "" {
		// Try Authorization header
		token = c.GetHeader("Authorization")
		if token != "" && len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	// TODO: Validate token if required
	// For now, we'll allow connections without strict token validation for testing

	// Serve WebSocket connection
	websocket.ServeWS(wh.hub, c.Writer, c.Request, sessionID, userID, videoID)
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
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebSocketAuthMiddleware struct {
	secretKey string
}

type WebSocketToken struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	VideoID   uuid.UUID `json:"video_id"`
	ExpiresAt int64     `json:"expires_at"`
	Scopes    []string  `json:"scopes"` // ["video_stream", "quality_control", "analytics"]
}

func NewWebSocketAuthMiddleware(secretKey string) *WebSocketAuthMiddleware {
	return &WebSocketAuthMiddleware{
		secretKey: secretKey,
	}
}

// GenerateWebSocketToken creates a signed token for WebSocket authentication
func (wsm *WebSocketAuthMiddleware) GenerateWebSocketToken(userID, sessionID, videoID uuid.UUID, scopes []string, duration time.Duration) (string, error) {
	token := WebSocketToken{
		UserID:    userID,
		SessionID: sessionID,
		VideoID:   videoID,
		ExpiresAt: time.Now().Add(duration).Unix(),
		Scopes:    scopes,
	}

	// Marshal to JSON
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token: %w", err)
	}

	// Base64 encode
	tokenB64 := base64.URLEncoding.EncodeToString(tokenJSON)

	// Generate signature
	signature := wsm.generateSignature(tokenB64)

	// Combine token and signature
	signedToken := fmt.Sprintf("%s.%s", tokenB64, signature)

	return signedToken, nil
}

// ValidateWebSocketToken validates a WebSocket token and returns the claims
func (wsm *WebSocketAuthMiddleware) ValidateWebSocketToken(tokenString string) (*WebSocketToken, error) {
	// Split token and signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	tokenB64 := parts[0]
	signature := parts[1]

	// Verify signature
	expectedSignature := wsm.generateSignature(tokenB64)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Decode token
	tokenJSON, err := base64.URLEncoding.DecodeString(tokenB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	// Unmarshal token
	var token WebSocketToken
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	// Check expiry
	if time.Now().Unix() > token.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}

// RequireWebSocketAuth middleware for WebSocket connections
func (wsm *WebSocketAuthMiddleware) RequireWebSocketAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get token from query parameter or header
		token := c.Query("ws_token")
		if token == "" {
			token = c.GetHeader("X-WS-Token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "WebSocket authentication token required",
				"code":  "WS_TOKEN_REQUIRED",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := wsm.ValidateWebSocketToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Invalid WebSocket token: %v", err),
				"code":  "WS_TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Check scope requirements
		requiredScope := c.Query("required_scope")
		if requiredScope != "" && !wsm.hasScope(claims.Scopes, requiredScope) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Insufficient scope: %s required", requiredScope),
				"code":  "WS_INSUFFICIENT_SCOPE",
			})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("ws_user_id", claims.UserID.String())
		c.Set("ws_session_id", claims.SessionID.String())
		c.Set("ws_video_id", claims.VideoID.String())
		c.Set("ws_scopes", claims.Scopes)
		c.Set("ws_expires_at", claims.ExpiresAt)

		c.Next()
	})
}

// OptionalWebSocketAuth middleware for optional WebSocket authentication
func (wsm *WebSocketAuthMiddleware) OptionalWebSocketAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get token from query parameter or header
		token := c.Query("ws_token")
		if token == "" {
			token = c.GetHeader("X-WS-Token")
		}

		if token == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Validate token if provided
		claims, err := wsm.ValidateWebSocketToken(token)
		if err != nil {
			// Invalid token, but we'll allow the connection with limited access
			c.Set("ws_auth_error", err.Error())
			c.Next()
			return
		}

		// Valid token - add claims to context
		c.Set("ws_user_id", claims.UserID.String())
		c.Set("ws_session_id", claims.SessionID.String())
		c.Set("ws_video_id", claims.VideoID.String())
		c.Set("ws_scopes", claims.Scopes)
		c.Set("ws_expires_at", claims.ExpiresAt)

		c.Next()
	})
}

// AuthenticateWebSocketUpgrade authenticates WebSocket upgrade requests
func (wsm *WebSocketAuthMiddleware) AuthenticateWebSocketUpgrade(w http.ResponseWriter, r *http.Request) (*WebSocketToken, error) {
	// Get token from query parameter or header
	token := r.URL.Query().Get("ws_token")
	if token == "" {
		token = r.Header.Get("X-WS-Token")
	}

	// Also check for Authorization header with Bearer token format
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:] // Remove "Bearer " prefix
		}
	}

	if token == "" {
		return nil, fmt.Errorf("no authentication token provided")
	}

	// Validate token
	claims, err := wsm.ValidateWebSocketToken(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return claims, nil
}

// ValidateWebSocketAccess validates specific WebSocket operations
func (wsm *WebSocketAuthMiddleware) ValidateWebSocketAccess(claims *WebSocketToken, operation string, resource string) error {
	// Check token expiry
	if time.Now().Unix() > claims.ExpiresAt {
		return fmt.Errorf("token expired")
	}

	// Check scope-based permissions
	switch operation {
	case "video_stream":
		if !wsm.hasScope(claims.Scopes, "video_stream") {
			return fmt.Errorf("video streaming scope required")
		}
	case "quality_control":
		if !wsm.hasScope(claims.Scopes, "quality_control") {
			return fmt.Errorf("quality control scope required")
		}
	case "analytics":
		if !wsm.hasScope(claims.Scopes, "analytics") {
			return fmt.Errorf("analytics scope required")
		}
	case "broadcast":
		if !wsm.hasScope(claims.Scopes, "broadcast") {
			return fmt.Errorf("broadcast scope required")
		}
	default:
		// Allow basic operations without specific scopes
	}

	return nil
}

// CreateSessionWebSocketToken creates a WebSocket token for a video session
func (wsm *WebSocketAuthMiddleware) CreateSessionWebSocketToken(userID, sessionID, videoID uuid.UUID, userRole string) (string, error) {
	// Determine scopes based on user role
	var scopes []string
	switch userRole {
	case "admin":
		scopes = []string{"video_stream", "quality_control", "analytics", "broadcast", "admin"}
	case "instructor":
		scopes = []string{"video_stream", "quality_control", "analytics", "broadcast"}
	default:
		scopes = []string{"video_stream"}
	}

	// Create token with 4-hour duration
	duration := 4 * time.Hour
	return wsm.GenerateWebSocketToken(userID, sessionID, videoID, scopes, duration)
}

// Helper methods

func (wsm *WebSocketAuthMiddleware) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(wsm.secretKey))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (wsm *WebSocketAuthMiddleware) hasScope(scopes []string, required string) bool {
	for _, scope := range scopes {
		if scope == required || scope == "admin" {
			return true
		}
	}
	return false
}

// WebSocket connection metadata for security logging
type WebSocketConnectionLog struct {
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	VideoID     uuid.UUID `json:"video_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	ConnectedAt time.Time `json:"connected_at"`
	Duration    int64     `json:"duration_seconds"`
	EventType   string    `json:"event_type"` // "connect", "disconnect", "error"
	ErrorReason string    `json:"error_reason,omitempty"`
}

// LogWebSocketConnection logs WebSocket connection events for security monitoring
func (wsm *WebSocketAuthMiddleware) LogWebSocketConnection(log *WebSocketConnectionLog) {
	// In production, this would send to a centralized logging system
	logJSON, _ := json.Marshal(log)
	fmt.Printf("WS_CONNECTION_LOG: %s\n", logJSON)
}

// RateLimitWebSocketConnections implements rate limiting for WebSocket connections
type WebSocketRateLimiter struct {
	connections map[string][]time.Time
	maxConns    int
	window      time.Duration
}

func NewWebSocketRateLimiter(maxConns int, window time.Duration) *WebSocketRateLimiter {
	return &WebSocketRateLimiter{
		connections: make(map[string][]time.Time),
		maxConns:    maxConns,
		window:      window,
	}
}

func (rl *WebSocketRateLimiter) AllowConnection(userID string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old connections
	if connections, exists := rl.connections[userID]; exists {
		var validConnections []time.Time
		for _, connTime := range connections {
			if connTime.After(cutoff) {
				validConnections = append(validConnections, connTime)
			}
		}
		rl.connections[userID] = validConnections
	}

	// Check if user exceeded rate limit
	if len(rl.connections[userID]) >= rl.maxConns {
		return false
	}

	// Add current connection
	rl.connections[userID] = append(rl.connections[userID], now)
	return true
}
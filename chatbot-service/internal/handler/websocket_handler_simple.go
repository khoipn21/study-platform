package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/repository"
	"chatbot-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var simpleUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin in development
	},
}

type SimpleWebSocketHandler struct {
	chatRepo  *repository.SimpleChatRepository
	aiService *service.AIService
	clients   map[*websocket.Conn]*SimpleClient
	mutex     sync.RWMutex
}

type SimpleClient struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Sessions map[uuid.UUID]bool // Track active sessions for this client
}

func NewSimpleWebSocketHandler(chatRepo *repository.SimpleChatRepository, aiService *service.AIService) *SimpleWebSocketHandler {
	return &SimpleWebSocketHandler{
		chatRepo:  chatRepo,
		aiService: aiService,
		clients:   make(map[*websocket.Conn]*SimpleClient),
	}
}

// extractUserIDFromToken extracts user_id from JWT token
// This is a simplified version - in production, use proper JWT library
func extractUserIDFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("Failed to decode token payload: %v", err)
		return ""
	}

	// Parse JSON
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		log.Printf("Failed to parse token claims: %v", err)
		return ""
	}

	// Extract user_id
	if userID, ok := claims["user_id"].(string); ok {
		return userID
	}
	if sub, ok := claims["sub"].(string); ok {
		return sub
	}

	return ""
}

func (h *SimpleWebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Try to get user ID from middleware first (X-User-ID header)
	userIDInterface, exists := c.Get("user_id")
	var userID uuid.UUID
	
	if !exists {
		// Fallback: Try to get from token query parameter or Authorization header
		token := c.Query("token")
		if token == "" {
			token = c.GetHeader("Authorization")
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}
		
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		
		// For demo, extract user_id from JWT claims
		// In production, properly validate JWT and extract claims
		userIDStr := extractUserIDFromToken(token)
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		userID = parsedID
	} else {
		var ok bool
		userID, ok = userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := simpleUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create and register client
	client := &SimpleClient{
		ID:       uuid.New(),
		UserID:   userID,
		Conn:     conn,
		Sessions: make(map[uuid.UUID]bool),
	}

	h.mutex.Lock()
	h.clients[conn] = client
	h.mutex.Unlock()

	// Clean up client when done
	defer func() {
		h.mutex.Lock()
		delete(h.clients, conn)
		h.mutex.Unlock()
	}()

	// Send welcome message
	welcomeMsg := map[string]interface{}{
		"type":    "connection",
		"status":  "connected",
		"message": "WebSocket connection established",
		"user_id": userID,
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	// Handle incoming messages
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle different message types
		h.handleMessage(client, msg)
	}
}

func (h *SimpleWebSocketHandler) handleMessage(client *SimpleClient, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		h.sendError(client.Conn, "Invalid message type")
		return
	}

	switch msgType {
	case "chat":
		h.handleChatMessage(client, msg)
	case "typing":
		h.handleTypingIndicator(client, msg)
	case "join_session":
		h.handleJoinSession(client, msg)
	case "leave_session":
		h.handleLeaveSession(client, msg)
	default:
		h.sendError(client.Conn, "Unknown message type")
	}
}

func (h *SimpleWebSocketHandler) handleChatMessage(client *SimpleClient, msg map[string]interface{}) {
	content, ok := msg["content"].(string)
	if !ok {
		h.sendError(client.Conn, "Invalid message content")
		return
	}

	sessionID := uuid.New()
	if sessionIDStr, exists := msg["session_id"].(string); exists {
		if parsed, err := uuid.Parse(sessionIDStr); err == nil {
			sessionID = parsed
		}
	}

	ctx := context.Background()

	// Save user message
	userMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Role:       model.RoleUser,
		Content:    content,
		TokensUsed: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessageWithUser(ctx, userMessage, client.UserID); err != nil {
		h.sendError(client.Conn, "Failed to save message")
		return
	}

	// Send user message confirmation
	h.sendMessage(client.Conn, map[string]interface{}{
		"type":       "message_saved",
		"message_id": userMessage.ID,
		"session_id": sessionID,
		"content":    content,
		"role":       "user",
		"created_at": userMessage.CreatedAt,
	})

	// Show typing indicator for AI
	h.sendMessage(client.Conn, map[string]interface{}{
		"type":       "typing",
		"session_id": sessionID,
		"sender":     "assistant",
		"typing":     true,
	})

	// Get recent messages for context
	messages, err := h.chatRepo.GetUserMessages(ctx, client.UserID, 10)
	if err != nil {
		h.sendError(client.Conn, "Failed to get message history")
		return
	}

	// Generate AI response
	aiResponse, err := h.aiService.GenerateResponse(ctx, messages, nil)
	if err != nil {
		h.sendError(client.Conn, "Failed to generate AI response")
		return
	}

	// Save AI message
	assistantMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Role:       model.RoleAssistant,
		Content:    aiResponse.Content,
		TokensUsed: aiResponse.TokensUsed,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessageWithUser(ctx, assistantMessage, client.UserID); err != nil {
		h.sendError(client.Conn, "Failed to save AI response")
		return
	}

	// Stop typing indicator
	h.sendMessage(client.Conn, map[string]interface{}{
		"type":       "typing",
		"session_id": sessionID,
		"sender":     "assistant",
		"typing":     false,
	})

	// Send AI response
	h.sendMessage(client.Conn, map[string]interface{}{
		"type":        "message",
		"message_id":  assistantMessage.ID,
		"session_id":  sessionID,
		"content":     aiResponse.Content,
		"role":        "assistant",
		"tokens_used": aiResponse.TokensUsed,
		"created_at":  assistantMessage.CreatedAt,
	})
}

func (h *SimpleWebSocketHandler) handleTypingIndicator(client *SimpleClient, msg map[string]interface{}) {
	sessionIDStr, ok := msg["session_id"].(string)
	if !ok {
		return
	}

	typing, ok := msg["typing"].(bool)
	if !ok {
		return
	}

	// Broadcast typing indicator to other clients in the same session
	h.broadcastToSession(sessionIDStr, map[string]interface{}{
		"type":       "typing",
		"session_id": sessionIDStr,
		"sender":     "user",
		"user_id":    client.UserID,
		"typing":     typing,
	}, client.Conn)
}

func (h *SimpleWebSocketHandler) handleJoinSession(client *SimpleClient, msg map[string]interface{}) {
	sessionIDStr, ok := msg["session_id"].(string)
	if !ok {
		h.sendError(client.Conn, "Invalid session ID")
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.sendError(client.Conn, "Invalid session ID format")
		return
	}

	client.Sessions[sessionID] = true

	h.sendMessage(client.Conn, map[string]interface{}{
		"type":       "session_joined",
		"session_id": sessionID,
		"status":     "success",
	})
}

func (h *SimpleWebSocketHandler) handleLeaveSession(client *SimpleClient, msg map[string]interface{}) {
	sessionIDStr, ok := msg["session_id"].(string)
	if !ok {
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return
	}

	delete(client.Sessions, sessionID)

	h.sendMessage(client.Conn, map[string]interface{}{
		"type":       "session_left",
		"session_id": sessionID,
		"status":     "success",
	})
}

func (h *SimpleWebSocketHandler) sendMessage(conn *websocket.Conn, msg map[string]interface{}) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *SimpleWebSocketHandler) sendError(conn *websocket.Conn, errorMsg string) {
	msg := map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	}
	h.sendMessage(conn, msg)
}

func (h *SimpleWebSocketHandler) broadcastToSession(sessionID string, msg map[string]interface{}, sender *websocket.Conn) {
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for conn, client := range h.clients {
		if conn != sender && client.Sessions[sessionUUID] {
			h.sendMessage(conn, msg)
		}
	}
}

// GetActiveClients returns the number of active WebSocket clients
func (h *SimpleWebSocketHandler) GetActiveClients() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// BroadcastMessage sends a message to all connected clients
func (h *SimpleWebSocketHandler) BroadcastMessage(msg map[string]interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for conn := range h.clients {
		h.sendMessage(conn, msg)
	}
}
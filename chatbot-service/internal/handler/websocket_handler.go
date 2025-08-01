package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/repository"
	"chatbot-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin in development
	},
}

type WebSocketHandler struct {
	chatRepo  *repository.ChatRepository
	aiService *service.AIService
	clients   map[*websocket.Conn]*Client
	mutex     sync.RWMutex
}

type Client struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Sessions map[uuid.UUID]bool // Track active sessions for this client
}

func NewWebSocketHandler(chatRepo *repository.ChatRepository, aiService *service.AIService) *WebSocketHandler {
	return &WebSocketHandler{
		chatRepo:  chatRepo,
		aiService: aiService,
		clients:   make(map[*websocket.Conn]*Client),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from middleware
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

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Create client
	client := &Client{
		ID:       uuid.New(),
		UserID:   userID,
		Conn:     conn,
		Sessions: make(map[uuid.UUID]bool),
	}

	// Register client
	h.mutex.Lock()
	h.clients[conn] = client
	h.mutex.Unlock()

	// Cleanup on disconnect
	defer func() {
		h.mutex.Lock()
		delete(h.clients, conn)
		h.mutex.Unlock()
	}()

	// Handle messages
	for {
		var wsMessage model.WebSocketMessage
		err := conn.ReadJSON(&wsMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		h.handleMessage(client, &wsMessage)
	}
}

func (h *WebSocketHandler) handleMessage(client *Client, wsMessage *model.WebSocketMessage) {
	ctx := context.Background()

	switch model.WebSocketMessageType(wsMessage.Type) {
	case model.WSMessageTypeChat:
		h.handleChatMessage(ctx, client, wsMessage)
	case model.WSMessageTypeSessionStart:
		h.handleSessionStart(ctx, client, wsMessage)
	case model.WSMessageTypeSessionEnd:
		h.handleSessionEnd(ctx, client, wsMessage)
	case model.WSMessageTypeTyping:
		h.handleTypingIndicator(ctx, client, wsMessage)
	default:
		h.sendError(client, "Unknown message type")
	}
}

func (h *WebSocketHandler) handleChatMessage(ctx context.Context, client *Client, wsMessage *model.WebSocketMessage) {
	// Parse chat request
	data, err := json.Marshal(wsMessage.Data)
	if err != nil {
		h.sendError(client, "Invalid message data")
		return
	}

	var chatReq model.ChatRequest
	if err := json.Unmarshal(data, &chatReq); err != nil {
		h.sendError(client, "Invalid chat request")
		return
	}

	// Validate session
	if chatReq.SessionID == nil {
		h.sendError(client, "Session ID is required")
		return
	}

	session, err := h.chatRepo.GetSessionByID(ctx, *chatReq.SessionID)
	if err != nil {
		h.sendError(client, "Session not found")
		return
	}

	if session.UserID != client.UserID {
		h.sendError(client, "Unauthorized access to session")
		return
	}

	// Mark client as active in this session
	client.Sessions[*chatReq.SessionID] = true

	// Send typing indicator
	h.sendTypingIndicator(client, *chatReq.SessionID, true)

	// Save user message
	userMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  *chatReq.SessionID,
		Role:       model.RoleUser,
		Content:    chatReq.Message,
		TokensUsed: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessage(ctx, userMessage); err != nil {
		h.sendError(client, "Failed to save message")
		return
	}

	// Get recent messages for context
	messages, err := h.chatRepo.GetRecentMessages(ctx, *chatReq.SessionID, 10)
	if err != nil {
		h.sendError(client, "Failed to get message history")
		return
	}

	// Build AI context
	var aiContext *service.ChatContext
	if chatReq.Context != nil {
		var contextData map[string]interface{}
		if err := json.Unmarshal([]byte(*chatReq.Context), &contextData); err == nil {
			aiContext = &service.ChatContext{}
			if courseName, ok := contextData["course_name"].(string); ok {
				aiContext.CourseName = courseName
			}
			if lectureName, ok := contextData["lecture_name"].(string); ok {
				aiContext.LectureName = lectureName
			}
		}
	}

	// Generate AI response (streaming)
	responseChan, errorChan := h.aiService.GenerateStreamResponse(ctx, messages, aiContext)

	var fullContent string
	var totalTokens int
	assistantMessageID := uuid.New()

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// Channel closed, save final message
				if fullContent != "" {
					assistantMessage := &model.ChatMessage{
						ID:         assistantMessageID,
						SessionID:  *chatReq.SessionID,
						Role:       model.RoleAssistant,
						Content:    fullContent,
						TokensUsed: totalTokens,
						CreatedAt:  time.Now(),
					}

					if err := h.chatRepo.CreateMessage(ctx, assistantMessage); err != nil {
						h.sendError(client, "Failed to save AI response")
						return
					}
				}

				// Stop typing indicator
				h.sendTypingIndicator(client, *chatReq.SessionID, false)
				return

			case err := <-errorChan:
				h.sendError(client, "AI service error: "+err.Error())
				h.sendTypingIndicator(client, *chatReq.SessionID, false)
				return

			default:
				if response != nil {
					fullContent += response.Content
					if response.TokensUsed > 0 {
						totalTokens = response.TokensUsed
					}

					// Send streaming response to client
					wsResponse := model.WebSocketMessage{
						Type:      string(model.WSMessageTypeResponse),
						SessionID: *chatReq.SessionID,
						Data: model.ChatResponse{
							SessionID:  *chatReq.SessionID,
							MessageID:  assistantMessageID,
							Role:       response.Role,
							Content:    response.Content,
							TokensUsed: response.TokensUsed,
							CreatedAt:  response.CreatedAt,
						},
					}

					if err := client.Conn.WriteJSON(wsResponse); err != nil {
						log.Printf("Failed to send WebSocket message: %v", err)
						return
					}
				}
			}
		}
	}
}

func (h *WebSocketHandler) handleSessionStart(ctx context.Context, client *Client, wsMessage *model.WebSocketMessage) {
	if wsMessage.SessionID != uuid.Nil {
		client.Sessions[wsMessage.SessionID] = true
	}
}

func (h *WebSocketHandler) handleSessionEnd(ctx context.Context, client *Client, wsMessage *model.WebSocketMessage) {
	if wsMessage.SessionID != uuid.Nil {
		delete(client.Sessions, wsMessage.SessionID)
	}
}

func (h *WebSocketHandler) handleTypingIndicator(ctx context.Context, client *Client, wsMessage *model.WebSocketMessage) {
	// Handle user typing indicator
	// In a real implementation, you might broadcast this to other users or moderators
}

func (h *WebSocketHandler) sendError(client *Client, message string) {
	wsMessage := model.WebSocketMessage{
		Type: string(model.WSMessageTypeError),
		Data: map[string]string{"error": message},
	}

	if err := client.Conn.WriteJSON(wsMessage); err != nil {
		log.Printf("Failed to send error message: %v", err)
	}
}

func (h *WebSocketHandler) sendTypingIndicator(client *Client, sessionID uuid.UUID, isTyping bool) {
	wsMessage := model.WebSocketMessage{
		Type:      string(model.WSMessageTypeTyping),
		SessionID: sessionID,
		Data: model.TypingIndicator{
			SessionID: sessionID,
			IsTyping:  isTyping,
		},
	}

	if err := client.Conn.WriteJSON(wsMessage); err != nil {
		log.Printf("Failed to send typing indicator: %v", err)
	}
}

func (h *WebSocketHandler) BroadcastToSession(sessionID uuid.UUID, message model.WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, client := range h.clients {
		if client.Sessions[sessionID] {
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Failed to broadcast message to client %v: %v", client.ID, err)
			}
		}
	}
}
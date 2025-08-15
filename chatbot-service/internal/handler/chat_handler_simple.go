package handler

import (
	"net/http"
	"strconv"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/repository"
	"chatbot-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SimpleChatHandler struct {
	chatRepo  *repository.SimpleChatRepository
	aiService *service.AIService
}

func NewSimpleChatHandler(chatRepo *repository.SimpleChatRepository, aiService *service.AIService) *SimpleChatHandler {
	return &SimpleChatHandler{
		chatRepo:  chatRepo,
		aiService: aiService,
	}
}

func (h *SimpleChatHandler) CreateSession(c *gin.Context) {
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

	var req struct {
		CourseID *uuid.UUID `json:"course_id,omitempty"`
		Title    string     `json:"title,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate title if not provided
	title := req.Title
	if title == "" {
		title = "New Chat Session"
	}

	// Create a fake session object since we don't have sessions table
	session := &model.ChatSession{
		ID:        uuid.New(),
		UserID:    userID,
		CourseID:  req.CourseID,
		Title:     title,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// No need to save to database since we don't have sessions table
	c.JSON(http.StatusCreated, session)
}

func (h *SimpleChatHandler) GetSession(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

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

	// Create a fake session
	session := &model.ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     "Chat Session",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Get messages if requested
	includeMessages := c.Query("include_messages") == "true"
	if includeMessages {
		messages, err := h.chatRepo.GetUserMessages(c.Request.Context(), userID, 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
			return
		}

		response := &model.ChatSessionResponse{
			ID:        session.ID,
			UserID:    session.UserID,
			CourseID:  session.CourseID,
			Title:     session.Title,
			IsActive:  session.IsActive,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		}

		// Convert to response messages
		for _, msg := range messages {
			response.Messages = append(response.Messages, *msg)
		}

		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SimpleChatHandler) GetUserSessions(c *gin.Context) {
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

	// Create a single fake session for this user
	sessions := []*model.ChatSession{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Title:     "Chat History",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"limit":    20,
		"offset":   0,
	})
}

func (h *SimpleChatHandler) UpdateSession(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

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

	var req struct {
		Title    *string `json:"title,omitempty"`
		IsActive *bool   `json:"is_active,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create fake updated session
	session := &model.ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     "Updated Chat Session",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Title != nil {
		session.Title = *req.Title
	}
	if req.IsActive != nil {
		session.IsActive = *req.IsActive
	}

	c.JSON(http.StatusOK, session)
}

func (h *SimpleChatHandler) DeleteSession(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	_, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

func (h *SimpleChatHandler) SendMessage(c *gin.Context) {
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

	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Generate session ID if not provided
	sessionID := uuid.New()
	if req.SessionID != nil {
		sessionID = *req.SessionID
	}

	// Save user message
	userMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Role:       model.RoleUser,
		Content:    req.Message,
		TokensUsed: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessageWithUser(c.Request.Context(), userMessage, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Get recent messages for context (last 10 messages from this user)
	messages, err := h.chatRepo.GetUserMessages(c.Request.Context(), userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message history"})
		return
	}

	// Generate AI response
	aiResponse, err := h.aiService.GenerateResponse(c.Request.Context(), messages, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI response", "details": err.Error()})
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

	if err := h.chatRepo.CreateMessageWithUser(c.Request.Context(), assistantMessage, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save AI response"})
		return
	}

	response := &model.ChatResponse{
		SessionID:  sessionID,
		MessageID:  assistantMessage.ID,
		Role:       model.RoleAssistant,
		Content:    aiResponse.Content,
		TokensUsed: aiResponse.TokensUsed,
		CreatedAt:  assistantMessage.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *SimpleChatHandler) GetMessages(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	_, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

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

	// Parse pagination parameters
	limit := 50
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get messages for this user
	messages, err := h.chatRepo.GetUserMessages(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}
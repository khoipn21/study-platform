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

type ChatHandler struct {
	chatRepo  *repository.ChatRepository
	aiService *service.AIService
}

func NewChatHandler(chatRepo *repository.ChatRepository, aiService *service.AIService) *ChatHandler {
	return &ChatHandler{
		chatRepo:  chatRepo,
		aiService: aiService,
	}
}

func (h *ChatHandler) CreateSession(c *gin.Context) {
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

	session := &model.ChatSession{
		ID:        uuid.New(),
		UserID:    userID,
		CourseID:  req.CourseID,
		Title:     title,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.chatRepo.CreateSession(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *ChatHandler) GetSession(c *gin.Context) {
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

	session, err := h.chatRepo.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get messages if requested
	includeMessages := c.Query("include_messages") == "true"
	if includeMessages {
		messages, err := h.chatRepo.GetSessionMessages(c.Request.Context(), sessionID, 100, 0)
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

func (h *ChatHandler) GetUserSessions(c *gin.Context) {
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
	limit := 20
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

	sessions, err := h.chatRepo.GetUserSessions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ChatHandler) UpdateSession(c *gin.Context) {
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

	// Get existing session
	session, err := h.chatRepo.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if req.Title != nil {
		session.Title = *req.Title
	}
	if req.IsActive != nil {
		session.IsActive = *req.IsActive
	}
	session.UpdatedAt = time.Now()

	if err := h.chatRepo.UpdateSession(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *ChatHandler) DeleteSession(c *gin.Context) {
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

	// Verify ownership
	session, err := h.chatRepo.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.chatRepo.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var session *model.ChatSession
	var err error

	// Create new session if sessionID is not provided
	if req.SessionID == nil {
		// Generate title from first message
		title, _ := h.aiService.GenerateSessionTitle(c.Request.Context(), req.Message)
		
		session = &model.ChatSession{
			ID:        uuid.New(),
			UserID:    userID,
			CourseID:  req.CourseID,
			Title:     title,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := h.chatRepo.CreateSession(c.Request.Context(), session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		req.SessionID = &session.ID
	} else {
		// Verify existing session
		session, err = h.chatRepo.GetSessionByID(c.Request.Context(), *req.SessionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		if session.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Save user message
	userMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  *req.SessionID,
		Role:       model.RoleUser,
		Content:    req.Message,
		TokensUsed: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessage(c.Request.Context(), userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Get recent messages for context
	messages, err := h.chatRepo.GetRecentMessages(c.Request.Context(), *req.SessionID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message history"})
		return
	}

	// Generate AI response
	aiResponse, err := h.aiService.GenerateResponse(c.Request.Context(), messages, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI response"})
		return
	}

	// Save AI message
	assistantMessage := &model.ChatMessage{
		ID:         uuid.New(),
		SessionID:  *req.SessionID,
		Role:       model.RoleAssistant,
		Content:    aiResponse.Content,
		TokensUsed: aiResponse.TokensUsed,
		CreatedAt:  time.Now(),
	}

	if err := h.chatRepo.CreateMessage(c.Request.Context(), assistantMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save AI response"})
		return
	}

	// Update session timestamp
	session.UpdatedAt = time.Now()
	h.chatRepo.UpdateSession(c.Request.Context(), session)

	response := &model.ChatResponse{
		SessionID:  *req.SessionID,
		MessageID:  assistantMessage.ID,
		Role:       model.RoleAssistant,
		Content:    aiResponse.Content,
		TokensUsed: aiResponse.TokensUsed,
		CreatedAt:  assistantMessage.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
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

	// Verify session ownership
	session, err := h.chatRepo.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
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

	messages, err := h.chatRepo.GetSessionMessages(c.Request.Context(), sessionID, limit, offset)
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
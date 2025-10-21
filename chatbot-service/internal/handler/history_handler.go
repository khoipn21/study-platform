package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"chatbot-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HistoryHandler struct {
	redisRepo *repository.RedisRepository
}

func NewHistoryHandler(redisRepo *repository.RedisRepository) *HistoryHandler {
	return &HistoryHandler{
		redisRepo: redisRepo,
	}
}

// ChatSessionInfo represents a chat session summary
type ChatSessionInfo struct {
	SessionID    string    `json:"session_id"`
	Title        string    `json:"title"`
	LastMessage  string    `json:"last_message"`
	MessageCount int       `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListUserSessions returns all chat sessions for a user from Redis
func (h *HistoryHandler) ListUserSessions(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	keys, err := h.redisRepo.ListChatSessions(ctx, userID)
	if err != nil {
		log.Printf("Error listing chat sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list chat sessions"})
		return
	}

	sessions := make([]ChatSessionInfo, 0, len(keys))
	for _, key := range keys {
		// Extract session ID from key: chat_history:{userid}:{sessionid}
		parts := strings.Split(key, ":")
		if len(parts) != 3 {
			continue
		}
		sessionID := parts[2]

		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			log.Printf("Invalid session UUID in key %s: %v", key, err)
			continue
		}

		// Get messages for this session
		messages, err := h.redisRepo.GetChatHistory(ctx, userID, sessionUUID)
		if err != nil {
			log.Printf("Error getting chat history for session %s: %v", sessionID, err)
			continue
		}

		if len(messages) == 0 {
			continue
		}

		// Get first user message as title
		title := "New Conversation"
		for _, msg := range messages {
			if msg.Role == "user" && len(msg.Content) > 0 {
				title = msg.Content
				if len(title) > 50 {
					title = title[:50] + "..."
				}
				break
			}
		}

		// Get last message preview
		lastMessage := ""
		if len(messages) > 0 {
			lastMsg := messages[len(messages)-1]
			lastMessage = lastMsg.Content
			if len(lastMessage) > 100 {
				lastMessage = lastMessage[:100] + "..."
			}
		}

		// Get timestamps
		createdAt := messages[0].CreatedAt
		updatedAt := messages[len(messages)-1].CreatedAt

		sessions = append(sessions, ChatSessionInfo{
			SessionID:    sessionID,
			Title:        title,
			LastMessage:  lastMessage,
			MessageCount: len(messages),
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	// Sort by updated time (most recent first)
	for i := 0; i < len(sessions)-1; i++ {
		for j := i + 1; j < len(sessions); j++ {
			if sessions[j].UpdatedAt.After(sessions[i].UpdatedAt) {
				sessions[i], sessions[j] = sessions[j], sessions[i]
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// GetSessionHistory returns all messages for a specific session
func (h *HistoryHandler) GetSessionHistory(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	ctx := context.Background()
	messages, err := h.redisRepo.GetChatHistory(ctx, userID, sessionID)
	if err != nil {
		log.Printf("Error getting chat history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionIDStr,
		"messages":   messages,
		"count":      len(messages),
	})
}

// DeleteSession deletes a chat session from Redis
func (h *HistoryHandler) DeleteSession(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	ctx := context.Background()
	err = h.redisRepo.DeleteChatSession(ctx, userID, sessionID)
	if err != nil {
		log.Printf("Error deleting chat session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Session %s deleted successfully", sessionIDStr),
	})
}

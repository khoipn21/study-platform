package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatSession struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CourseID    *uuid.UUID `json:"course_id,omitempty" db:"course_id"`
	Title       string    `json:"title" db:"title"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ChatMessage struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	SessionID     uuid.UUID    `json:"session_id" db:"session_id"`
	Role          MessageRole  `json:"role" db:"role"`
	Content       string       `json:"content" db:"content"`
	TokensUsed    int          `json:"tokens_used" db:"tokens_used"`
	ResponseTime  *int         `json:"response_time,omitempty" db:"response_time"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
}

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

type ChatRequest struct {
	SessionID *uuid.UUID `json:"session_id,omitempty"`
	CourseID  *uuid.UUID `json:"course_id,omitempty"`
	Message   string     `json:"message" binding:"required"`
	Context   *string    `json:"context,omitempty"`
}

type ChatResponse struct {
	SessionID uuid.UUID    `json:"session_id"`
	MessageID uuid.UUID    `json:"message_id"`
	Role      MessageRole  `json:"role"`
	Content   string       `json:"content"`
	TokensUsed int         `json:"tokens_used"`
	CreatedAt time.Time    `json:"created_at"`
}

type ChatSessionResponse struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	CourseID  *uuid.UUID    `json:"course_id,omitempty"`
	Title     string        `json:"title"`
	IsActive  bool          `json:"is_active"`
	Messages  []ChatMessage `json:"messages,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	SessionID uuid.UUID   `json:"session_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type WebSocketMessageType string

const (
	WSMessageTypeChat         WebSocketMessageType = "chat"
	WSMessageTypeResponse     WebSocketMessageType = "response"
	WSMessageTypeError        WebSocketMessageType = "error"
	WSMessageTypeTyping       WebSocketMessageType = "typing"
	WSMessageTypeSessionStart WebSocketMessageType = "session_start"
	WSMessageTypeSessionEnd   WebSocketMessageType = "session_end"
)

type TypingIndicator struct {
	SessionID uuid.UUID `json:"session_id"`
	IsTyping  bool      `json:"is_typing"`
}
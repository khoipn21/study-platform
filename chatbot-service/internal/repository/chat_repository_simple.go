package repository

import (
	"context"
	"database/sql"
	"time"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// SimpleChatRepository works with the existing chat_history table schema
type SimpleChatRepository struct {
	db *sql.DB
}

func NewSimpleChatRepository(db *sql.DB) *SimpleChatRepository {
	return &SimpleChatRepository{db: db}
}

// CreateSession is a no-op since we don't have sessions table - sessions are implicit
func (r *SimpleChatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	// No sessions table, so this is a no-op
	return nil
}

// GetSessionByID creates a fake session for compatibility
func (r *SimpleChatRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.ChatSession, error) {
	// Create a fake session object for compatibility
	return &model.ChatSession{
		ID:        sessionID,
		UserID:    uuid.New(), // This will be overridden by the handler
		Title:     "Chat Session",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetUserSessions returns empty list since we don't track sessions
func (r *SimpleChatRepository) GetUserSessions(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*model.ChatSession, error) {
	// Return a single default session for this user
	return []*model.ChatSession{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Title:     "Chat History",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// UpdateSession is a no-op since we don't have sessions table
func (r *SimpleChatRepository) UpdateSession(ctx context.Context, session *model.ChatSession) error {
	return nil
}

// DeleteSession is a no-op since we don't have sessions table
func (r *SimpleChatRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return nil
}

// CreateMessage saves a message to chat_history table
func (r *SimpleChatRepository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	query := `
		INSERT INTO chat_history (id, user_id, message, is_user, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	// Determine if it's a user message
	isUser := message.Role == model.RoleUser

	// We need to get the user_id from somewhere - let's use a dummy user for now
	// In a real implementation, this would come from the session context
	userID := uuid.New() // This should be passed as parameter or retrieved from context

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		userID,
		message.Content,
		isUser,
		message.CreatedAt,
	)

	return err
}

// CreateMessageWithUser saves a message to chat_history table with specific user ID
func (r *SimpleChatRepository) CreateMessageWithUser(ctx context.Context, message *model.ChatMessage, userID uuid.UUID) error {
	query := `
		INSERT INTO chat_history (user_id, message, is_user)
		VALUES ($1, $2, $3)`

	// Determine if it's a user message
	isUser := message.Role == model.RoleUser

	_, err := r.db.ExecContext(ctx, query,
		userID,
		message.Content,
		isUser,
	)

	return err
}

// GetSessionMessages retrieves messages for a user (treating user as session)
func (r *SimpleChatRepository) GetSessionMessages(ctx context.Context, sessionID uuid.UUID, limit int, offset int) ([]*model.ChatMessage, error) {
	// For simplicity, get recent messages for any user
	query := `
		SELECT id, user_id, message, is_user, created_at
		FROM chat_history
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var id uuid.UUID
		var userID uuid.UUID
		var content string
		var isUser bool
		var createdAt time.Time

		err := rows.Scan(&id, &userID, &content, &isUser, &createdAt)
		if err != nil {
			return nil, err
		}

		role := model.RoleAssistant
		if isUser {
			role = model.RoleUser
		}

		message := &model.ChatMessage{
			ID:         id,
			SessionID:  sessionID, // Use the provided sessionID
			Role:       role,
			Content:    content,
			TokensUsed: 0, // Not tracked in simple schema
			CreatedAt:  createdAt,
		}

		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetRecentMessages retrieves recent messages for context
func (r *SimpleChatRepository) GetRecentMessages(ctx context.Context, sessionID uuid.UUID, limit int) ([]*model.ChatMessage, error) {
	return r.GetSessionMessages(ctx, sessionID, limit, 0)
}

// GetUserMessages retrieves messages for a specific user
func (r *SimpleChatRepository) GetUserMessages(ctx context.Context, userID uuid.UUID, limit int) ([]*model.ChatMessage, error) {
	query := `
		SELECT id, user_id, message, is_user, created_at
		FROM chat_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var id uuid.UUID
		var userID uuid.UUID
		var content string
		var isUser bool
		var createdAt time.Time

		err := rows.Scan(&id, &userID, &content, &isUser, &createdAt)
		if err != nil {
			return nil, err
		}

		role := model.RoleAssistant
		if isUser {
			role = model.RoleUser
		}

		message := &model.ChatMessage{
			ID:         id,
			SessionID:  uuid.New(), // Generate a dummy session ID
			Role:       role,
			Content:    content,
			TokensUsed: 0, // Not tracked in simple schema
			CreatedAt:  createdAt,
		}

		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetUserTokenUsage returns 0 since we don't track tokens in simple schema
func (r *SimpleChatRepository) GetUserTokenUsage(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error) {
	return 0, nil
}

// CleanupOldSessions deletes old chat_history entries
func (r *SimpleChatRepository) CleanupOldSessions(ctx context.Context, olderThan time.Time) error {
	query := `DELETE FROM chat_history WHERE created_at < $1`
	_, err := r.db.ExecContext(ctx, query, olderThan)
	return err
}
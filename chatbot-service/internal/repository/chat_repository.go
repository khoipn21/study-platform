package repository

import (
	"context"
	"database/sql"
	"time"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (id, user_id, course_id, title, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.CourseID,
		session.Title,
		session.IsActive,
		session.CreatedAt,
		session.UpdatedAt,
	)

	return err
}

func (r *ChatRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.ChatSession, error) {
	query := `
		SELECT id, user_id, course_id, title, is_active, created_at, updated_at
		FROM chat_sessions
		WHERE id = $1 AND deleted_at IS NULL`

	session := &model.ChatSession{}
	var courseID sql.NullString

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&courseID,
		&session.Title,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if courseID.Valid {
		courseUUID, err := uuid.Parse(courseID.String)
		if err == nil {
			session.CourseID = &courseUUID
		}
	}

	return session, nil
}

func (r *ChatRepository) GetUserSessions(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*model.ChatSession, error) {
	query := `
		SELECT id, user_id, course_id, title, is_active, created_at, updated_at
		FROM chat_sessions
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.ChatSession
	for rows.Next() {
		session := &model.ChatSession{}
		var courseID sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&courseID,
			&session.Title,
			&session.IsActive,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if courseID.Valid {
			courseUUID, err := uuid.Parse(courseID.String)
			if err == nil {
				session.CourseID = &courseUUID
			}
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *ChatRepository) UpdateSession(ctx context.Context, session *model.ChatSession) error {
	query := `
		UPDATE chat_sessions
		SET title = $2, is_active = $3, updated_at = $4
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.Title,
		session.IsActive,
		session.UpdatedAt,
	)

	return err
}

func (r *ChatRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		UPDATE chat_sessions
		SET deleted_at = $2
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, sessionID, time.Now())
	return err
}

func (r *ChatRepository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (id, session_id, role, content, tokens_used, response_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.SessionID,
		message.Role,
		message.Content,
		message.TokensUsed,
		message.ResponseTime,
		message.CreatedAt,
	)

	return err
}

func (r *ChatRepository) GetSessionMessages(ctx context.Context, sessionID uuid.UUID, limit int, offset int) ([]*model.ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, tokens_used, response_time, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		message := &model.ChatMessage{}
		var responseTime sql.NullInt32

		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.TokensUsed,
			&responseTime,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if responseTime.Valid {
			responseTimeInt := int(responseTime.Int32)
			message.ResponseTime = &responseTimeInt
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *ChatRepository) GetRecentMessages(ctx context.Context, sessionID uuid.UUID, limit int) ([]*model.ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, tokens_used, response_time, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		message := &model.ChatMessage{}
		var responseTime sql.NullInt32

		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.TokensUsed,
			&responseTime,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if responseTime.Valid {
			responseTimeInt := int(responseTime.Int32)
			message.ResponseTime = &responseTimeInt
		}

		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) GetUserTokenUsage(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(cm.tokens_used), 0)
		FROM chat_messages cm
		JOIN chat_sessions cs ON cm.session_id = cs.id
		WHERE cs.user_id = $1 AND cm.created_at BETWEEN $2 AND $3`

	var totalTokens int
	err := r.db.QueryRowContext(ctx, query, userID, from, to).Scan(&totalTokens)
	return totalTokens, err
}

func (r *ChatRepository) CleanupOldSessions(ctx context.Context, olderThan time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete old messages first (foreign key constraint)
	messageQuery := `
		DELETE FROM chat_messages
		WHERE session_id IN (
			SELECT id FROM chat_sessions
			WHERE updated_at < $1 AND is_active = false
		)`

	_, err = tx.ExecContext(ctx, messageQuery, olderThan)
	if err != nil {
		return err
	}

	// Delete old sessions
	sessionQuery := `
		DELETE FROM chat_sessions
		WHERE updated_at < $1 AND is_active = false`

	_, err = tx.ExecContext(ctx, sessionQuery, olderThan)
	if err != nil {
		return err
	}

	return tx.Commit()
}
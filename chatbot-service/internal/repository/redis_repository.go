package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

// StoreChatHistory stores chat history in Redis with the key format chat_history:{userid}:{chatid}
func (r *RedisRepository) StoreChatHistory(ctx context.Context, accountID uuid.UUID, chatID uuid.UUID, messages []model.ChatMessage) error {
	key := fmt.Sprintf("chat_history:%s:%s", accountID.String(), chatID.String())
	
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	// Store as JSON string with 30 days expiration
	err = r.client.Set(ctx, key, data, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store chat history: %w", err)
	}

	return nil
}

// GetChatHistory retrieves chat history from Redis
func (r *RedisRepository) GetChatHistory(ctx context.Context, accountID uuid.UUID, chatID uuid.UUID) ([]model.ChatMessage, error) {
	key := fmt.Sprintf("chat_history:%s:%s", accountID.String(), chatID.String())
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return []model.ChatMessage{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	var messages []model.ChatMessage
	err = json.Unmarshal([]byte(data), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	return messages, nil
}

// AppendMessage adds a new message to existing chat history
func (r *RedisRepository) AppendMessage(ctx context.Context, accountID uuid.UUID, chatID uuid.UUID, message model.ChatMessage) error {
	messages, err := r.GetChatHistory(ctx, accountID, chatID)
	if err != nil {
		return err
	}

	messages = append(messages, message)
	return r.StoreChatHistory(ctx, accountID, chatID, messages)
}

// ListChatSessions returns all chat sessions for an account
func (r *RedisRepository) ListChatSessions(ctx context.Context, accountID uuid.UUID) ([]string, error) {
	pattern := fmt.Sprintf("chat_history:%s:*", accountID.String())
	
	var keys []string
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan keys: %w", err)
	}

	return keys, nil
}

// DeleteChatSession deletes a chat session
func (r *RedisRepository) DeleteChatSession(ctx context.Context, accountID uuid.UUID, chatID uuid.UUID) error {
	key := fmt.Sprintf("chat_history:%s:%s", accountID.String(), chatID.String())
	return r.client.Del(ctx, key).Err()
}

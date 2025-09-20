package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"video-service/internal/config"
	"video-service/internal/model"
)

type RedisClient struct {
	client *redis.Client
	config *config.Config
}

type Publisher interface {
	PublishNetworkStatus(ctx context.Context, sessionID string, status *model.WSNetworkStatus) error
	PublishQualityChange(ctx context.Context, sessionID string, change *model.WSQualityChange) error
	PublishAnalyticsEvent(ctx context.Context, event *model.WSAnalyticsEvent) error
	PublishHeartbeat(ctx context.Context, sessionID string, data map[string]interface{}) error
}

type Subscriber interface {
	SubscribeNetworkStatus(ctx context.Context, sessionID string, callback func(*model.WSNetworkStatus)) error
	SubscribeQualityChange(ctx context.Context, sessionID string, callback func(*model.WSQualityChange)) error
	SubscribeAnalytics(ctx context.Context, callback func(*model.WSAnalyticsEvent)) error
	SubscribeHeartbeat(ctx context.Context, callback func(string, map[string]interface{})) error
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: cfg.RedisPoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: rdb,
		config: cfg,
	}, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Publisher implementation
func (r *RedisClient) PublishNetworkStatus(ctx context.Context, sessionID string, status *model.WSNetworkStatus) error {
	channel := fmt.Sprintf("network_status:%s", sessionID)
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal network status: %w", err)
	}

	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisClient) PublishQualityChange(ctx context.Context, sessionID string, change *model.WSQualityChange) error {
	channel := fmt.Sprintf("quality_change:%s", sessionID)
	data, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("failed to marshal quality change: %w", err)
	}

	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisClient) PublishAnalyticsEvent(ctx context.Context, event *model.WSAnalyticsEvent) error {
	channel := "video_analytics"
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal analytics event: %w", err)
	}

	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisClient) PublishHeartbeat(ctx context.Context, sessionID string, data map[string]interface{}) error {
	channel := "session_heartbeat"
	payload := map[string]interface{}{
		"session_id": sessionID,
		"data":       data,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat: %w", err)
	}

	return r.client.Publish(ctx, channel, jsonData).Err()
}

// Subscriber implementation
func (r *RedisClient) SubscribeNetworkStatus(ctx context.Context, sessionID string, callback func(*model.WSNetworkStatus)) error {
	channel := fmt.Sprintf("network_status:%s", sessionID)
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			var status model.WSNetworkStatus
			if err := json.Unmarshal([]byte(msg.Payload), &status); err != nil {
				continue
			}
			callback(&status)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *RedisClient) SubscribeQualityChange(ctx context.Context, sessionID string, callback func(*model.WSQualityChange)) error {
	channel := fmt.Sprintf("quality_change:%s", sessionID)
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			var change model.WSQualityChange
			if err := json.Unmarshal([]byte(msg.Payload), &change); err != nil {
				continue
			}
			callback(&change)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *RedisClient) SubscribeAnalytics(ctx context.Context, callback func(*model.WSAnalyticsEvent)) error {
	channel := "video_analytics"
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			var event model.WSAnalyticsEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				continue
			}
			callback(&event)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *RedisClient) SubscribeHeartbeat(ctx context.Context, callback func(string, map[string]interface{})) error {
	channel := "session_heartbeat"
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				continue
			}

			sessionID, ok := payload["session_id"].(string)
			if !ok {
				continue
			}

			data, ok := payload["data"].(map[string]interface{})
			if !ok {
				continue
			}

			callback(sessionID, data)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Cache methods for storing session data
func (r *RedisClient) SetSessionData(ctx context.Context, sessionID string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisClient) GetSessionData(ctx context.Context, sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Del(ctx, key).Err()
}

// Network metrics caching
func (r *RedisClient) CacheNetworkMetrics(ctx context.Context, sessionID string, metrics *model.NetworkMetrics, expiration time.Duration) error {
	key := fmt.Sprintf("metrics:%s", sessionID)
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal network metrics: %w", err)
	}

	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisClient) GetCachedNetworkMetrics(ctx context.Context, sessionID string) (*model.NetworkMetrics, error) {
	key := fmt.Sprintf("metrics:%s", sessionID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var metrics model.NetworkMetrics
	if err := json.Unmarshal([]byte(data), &metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}

// Quality recommendation caching
func (r *RedisClient) CacheQualityRecommendation(ctx context.Context, sessionID, quality string, expiration time.Duration) error {
	key := fmt.Sprintf("quality_rec:%s", sessionID)
	return r.client.Set(ctx, key, quality, expiration).Err()
}

func (r *RedisClient) GetCachedQualityRecommendation(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("quality_rec:%s", sessionID)
	return r.client.Get(ctx, key).Result()
}

// Active sessions management
func (r *RedisClient) AddActiveSession(ctx context.Context, sessionID, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("active_sessions:%s", userID)
	return r.client.SAdd(ctx, key, sessionID).Err()
}

func (r *RedisClient) RemoveActiveSession(ctx context.Context, sessionID, userID string) error {
	key := fmt.Sprintf("active_sessions:%s", userID)
	return r.client.SRem(ctx, key, sessionID).Err()
}

func (r *RedisClient) GetActiveSessions(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf("active_sessions:%s", userID)
	return r.client.SMembers(ctx, key).Result()
}

// General key-value operations
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
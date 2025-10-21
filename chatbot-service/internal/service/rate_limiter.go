package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client        *redis.Client
	dailyLimit    int
	limiterWindow time.Duration
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client:        client,
		dailyLimit:    10, // 10 prompts per day
		limiterWindow: 24 * time.Hour,
	}
}

// CheckLimit checks if the user has exceeded their daily limit
func (rl *RateLimiter) CheckLimit(ctx context.Context, accountID uuid.UUID) (bool, int, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", accountID.String(), time.Now().Format("2006-01-02"))
	
	count, err := rl.client.Get(ctx, key).Int()
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return false, 0, fmt.Errorf("failed to get rate limit: %w", err)
	}

	remaining := rl.dailyLimit - count
	if remaining <= 0 {
		return false, 0, nil
	}

	return true, remaining, nil
}

// IncrementUsage increments the usage count for a user
func (rl *RateLimiter) IncrementUsage(ctx context.Context, accountID uuid.UUID) error {
	key := fmt.Sprintf("rate_limit:%s:%s", accountID.String(), time.Now().Format("2006-01-02"))
	
	pipe := rl.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.limiterWindow)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}

	return nil
}

// GetUsage returns the current usage count for a user
func (rl *RateLimiter) GetUsage(ctx context.Context, accountID uuid.UUID) (int, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", accountID.String(), time.Now().Format("2006-01-02"))
	
	count, err := rl.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get usage: %w", err)
	}

	return count, nil
}

// GetRemaining returns the remaining prompts for today
func (rl *RateLimiter) GetRemaining(ctx context.Context, accountID uuid.UUID) (int, error) {
	usage, err := rl.GetUsage(ctx, accountID)
	if err != nil {
		return 0, err
	}

	remaining := rl.dailyLimit - usage
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// ResetLimit resets the limit for a user (admin function)
func (rl *RateLimiter) ResetLimit(ctx context.Context, accountID uuid.UUID) error {
	key := fmt.Sprintf("rate_limit:%s:%s", accountID.String(), time.Now().Format("2006-01-02"))
	return rl.client.Del(ctx, key).Err()
}

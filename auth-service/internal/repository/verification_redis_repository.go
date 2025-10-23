package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type VerificationRedisRepository struct {
	client *redis.Client
}

func NewVerificationRedisRepository(client *redis.Client) *VerificationRedisRepository {
	return &VerificationRedisRepository{client: client}
}

// Redis Key Patterns:
// otp:{user_id}               -> "123456" (6-digit code, TTL: 15 min)
// otp_attempts:{user_id}      -> "3" (attempt counter, TTL: 15 min)
// otp_token:{token}           -> user_id (for email link verification, TTL: 24h)
// resend_limit:{email}        -> "2" (resend counter, TTL: 1 hour)

// GenerateOTP creates a 6-digit numeric code
func (r *VerificationRedisRepository) GenerateOTP() string {
	otp := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += fmt.Sprintf("%d", n.Int64())
	}
	return otp
}

// GenerateToken creates a URL-safe token
func (r *VerificationRedisRepository) GenerateToken() string {
	b := make([]byte, 48)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// StoreOTP stores OTP with 15-minute expiry
func (r *VerificationRedisRepository) StoreOTP(ctx context.Context, userID uuid.UUID, otp string) error {
	key := fmt.Sprintf("otp:%s", userID.String())
	return r.client.Set(ctx, key, otp, 15*time.Minute).Err()
}

// GetOTP retrieves stored OTP
func (r *VerificationRedisRepository) GetOTP(ctx context.Context, userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("otp:%s", userID.String())
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no OTP found or expired")
	}
	return result, err
}

// DeleteOTP removes OTP after successful verification
func (r *VerificationRedisRepository) DeleteOTP(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("otp:%s", userID.String())
	return r.client.Del(ctx, key).Err()
}

// IncrementAttempts increments and returns attempt count
func (r *VerificationRedisRepository) IncrementAttempts(ctx context.Context, userID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s", userID.String())

	// Increment
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiry on first attempt (only if not already set)
	if count == 1 {
		r.client.Expire(ctx, key, 15*time.Minute)
	}

	return count, nil
}

// GetAttempts retrieves current attempt count
func (r *VerificationRedisRepository) GetAttempts(ctx context.Context, userID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s", userID.String())
	result, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

// StoreVerificationToken stores token mapping to user_id (for email link clicks)
func (r *VerificationRedisRepository) StoreVerificationToken(ctx context.Context, token string, userID uuid.UUID) error {
	key := fmt.Sprintf("otp_token:%s", token)
	return r.client.Set(ctx, key, userID.String(), 24*time.Hour).Err()
}

// GetUserIDFromToken retrieves user_id from token
func (r *VerificationRedisRepository) GetUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	key := fmt.Sprintf("otp_token:%s", token)
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return uuid.Nil, fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token data")
	}

	return userID, nil
}

// DeleteVerificationToken removes token after use (one-time use)
func (r *VerificationRedisRepository) DeleteVerificationToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("otp_token:%s", token)
	return r.client.Del(ctx, key).Err()
}

// CheckResendLimit checks if user can resend (max 3 per hour)
func (r *VerificationRedisRepository) CheckResendLimit(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("resend_limit:%s", email)
	count, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return true, nil // No limit hit yet
	}
	if err != nil {
		return false, err
	}

	return count < 3, nil
}

// IncrementResendLimit increments resend counter
func (r *VerificationRedisRepository) IncrementResendLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("resend_limit:%s", email)

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// Set 1-hour expiry on first resend
	if count == 1 {
		r.client.Expire(ctx, key, 1*time.Hour)
	}

	return nil
}

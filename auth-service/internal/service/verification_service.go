package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/auth-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type VerificationService struct {
	redisRepo    *repository.VerificationRedisRepository
	userRepo     *repository.UserRepository
	emailService *EmailService
	logger       logger.Logger
}

func NewVerificationService(
	redisRepo *repository.VerificationRedisRepository,
	userRepo *repository.UserRepository,
	emailService *EmailService,
	logger logger.Logger,
) *VerificationService {
	return &VerificationService{
		redisRepo:    redisRepo,
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// SendVerification generates OTP and sends verification email
func (s *VerificationService) SendVerification(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Generate OTP and token
	otp := s.redisRepo.GenerateOTP()
	token := s.redisRepo.GenerateToken()

	// Store in Redis
	if err := s.redisRepo.StoreOTP(ctx, userID, otp); err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	if err := s.redisRepo.StoreVerificationToken(ctx, token, userID); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Send email
	if err := s.emailService.SendVerificationEmail(user.Email, user.Username, otp, token); err != nil {
		s.logger.Errorf("Failed to send verification email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Infof("Verification email sent to: %s", user.Email)
	return nil
}

// VerifyOTP verifies the OTP code
func (s *VerificationService) VerifyOTP(ctx context.Context, userID uuid.UUID, inputOTP string) error {
	// Check attempts
	attempts, err := s.redisRepo.GetAttempts(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check attempts: %w", err)
	}

	if attempts >= 5 {
		return fmt.Errorf("too_many_attempts")
	}

	// Get stored OTP
	storedOTP, err := s.redisRepo.GetOTP(ctx, userID)
	if err != nil {
		return fmt.Errorf("otp_expired")
	}

	// Verify OTP
	if storedOTP != inputOTP {
		s.redisRepo.IncrementAttempts(ctx, userID)
		return fmt.Errorf("invalid_otp")
	}

	// Success - delete OTP and mark user as verified
	s.redisRepo.DeleteOTP(ctx, userID)
	if err := s.userRepo.MarkEmailAsVerified(userID); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Infof("Email verified for user: %s", userID.String())
	return nil
}

// VerifyToken verifies via email link click
func (s *VerificationService) VerifyToken(ctx context.Context, token string) (*uuid.UUID, error) {
	// Get user ID from token
	userID, err := s.redisRepo.GetUserIDFromToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token_invalid")
	}

	// Delete token (one-time use)
	s.redisRepo.DeleteVerificationToken(ctx, token)

	// Mark user as verified
	if err := s.userRepo.MarkEmailAsVerified(userID); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Infof("Email verified via token for user: %s", userID.String())
	return &userID, nil
}

// ResendVerification resends verification email with rate limiting
func (s *VerificationService) ResendVerification(ctx context.Context, email string) error {
	// Check rate limit
	canResend, err := s.redisRepo.CheckResendLimit(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if !canResend {
		return fmt.Errorf("rate_limit")
	}

	// Get user
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.IsEmailVerified {
		return fmt.Errorf("already_verified")
	}

	// Increment resend counter
	if err := s.redisRepo.IncrementResendLimit(ctx, email); err != nil {
		return fmt.Errorf("failed to update rate limit: %w", err)
	}

	// Send new verification
	return s.SendVerification(ctx, user.ID)
}

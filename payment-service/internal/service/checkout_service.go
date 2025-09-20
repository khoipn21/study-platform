package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"
)

// CheckoutOptions contains options for creating checkout sessions
type CheckoutOptions struct {
	SuccessURL    string                 `json:"success_url,omitempty"`
	CancelURL     string                 `json:"cancel_url,omitempty"`
	CustomerEmail string                 `json:"customer_email,omitempty"`
	CustomerName  string                 `json:"customer_name,omitempty"`
	CustomData    map[string]interface{} `json:"custom_data,omitempty"`
}

type CheckoutService struct {
	transactionRepo    *repository.TransactionRepository
	courseCheckoutSvc  *CourseCheckoutService
}

func NewCheckoutService(transactionRepo *repository.TransactionRepository, courseCheckoutSvc *CourseCheckoutService) *CheckoutService {
	return &CheckoutService{
		transactionRepo:   transactionRepo,
		courseCheckoutSvc: courseCheckoutSvc,
	}
}

// CreateCheckoutSession creates a new checkout session
func (s *CheckoutService) CreateCheckoutSession(ctx context.Context, userID, courseID string, amount float64, currency string) (*model.CheckoutSession, error) {
	session := &model.CheckoutSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		CourseID:  courseID,
		Amount:    amount,
		Currency:  currency,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute), // 30-minute expiry
	}

	return session, nil
}

// ValidateCheckout validates checkout session
func (s *CheckoutService) ValidateCheckout(ctx context.Context, sessionID string) (*model.CheckoutSession, error) {
	// In a real implementation, this would validate against stored sessions
	return nil, fmt.Errorf("checkout session not found")
}

// CompleteCheckout completes a checkout session
func (s *CheckoutService) CompleteCheckout(ctx context.Context, sessionID string, paymentData map[string]interface{}) error {
	// In a real implementation, this would update the session status
	return nil
}

// CreateCourseCheckout creates a checkout session for a course purchase
func (s *CheckoutService) CreateCourseCheckout(ctx context.Context, userID, courseID string, options *CheckoutOptions) (*CourseCheckoutResponse, error) {
	if s.courseCheckoutSvc == nil {
		return nil, fmt.Errorf("course checkout service not initialized")
	}

	req := &CourseCheckoutRequest{
		UserID:        userID,
		CourseID:      courseID,
		SuccessURL:    options.SuccessURL,
		CancelURL:     options.CancelURL,
		CustomerEmail: options.CustomerEmail,
		CustomerName:  options.CustomerName,
	}

	return s.courseCheckoutSvc.CreateCourseCheckout(ctx, req)
}
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/lemonsqueezy"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"
)

// CourseCheckoutService handles LemonSqueezy checkout creation for courses
type CourseCheckoutService struct {
	lemonSqueezyClient *lemonsqueezy.Client
	transactionRepo    *repository.TransactionRepository
	enrollmentRepo     *repository.EnrollmentRepository
	courseRepo         *repository.CourseRepository
	auditRepo          *repository.AuditRepository
}

type CourseCheckoutRequest struct {
	CourseID      string `json:"course_id" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	SuccessURL    string `json:"success_url,omitempty"`
	CancelURL     string `json:"cancel_url,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
}

type CourseCheckoutResponse struct {
	CheckoutURL    string    `json:"checkout_url"`
	CheckoutID     string    `json:"checkout_id"`
	TransactionID  string    `json:"transaction_id"`
	CourseTitle    string    `json:"course_title"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewCourseCheckoutService(
	lemonSqueezyClient *lemonsqueezy.Client,
	transactionRepo *repository.TransactionRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	courseRepo *repository.CourseRepository,
	auditRepo *repository.AuditRepository,
) *CourseCheckoutService {
	return &CourseCheckoutService{
		lemonSqueezyClient: lemonSqueezyClient,
		transactionRepo:    transactionRepo,
		enrollmentRepo:     enrollmentRepo,
		courseRepo:         courseRepo,
		auditRepo:          auditRepo,
	}
}

// CreateCourseCheckout creates a LemonSqueezy checkout session for a course purchase
func (s *CourseCheckoutService) CreateCourseCheckout(ctx context.Context, req *CourseCheckoutRequest) (*CourseCheckoutResponse, error) {
	// 1. Validate and get course information
	course, err := s.courseRepo.GetByID(ctx, req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// 2. Validate course is available for purchase
	if course.Status != "published" {
		return nil, fmt.Errorf("course is not available for purchase")
	}

	if !course.IsPaid {
		return nil, fmt.Errorf("course is free and does not require payment")
	}

	if course.LemonSqueezyProductID == nil || course.LemonSqueezyVariantID == nil {
		return nil, fmt.Errorf("course is not configured for LemonSqueezy payments")
	}

	// 3. Check if user is already enrolled
	existingEnrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, req.UserID, req.CourseID)
	if err == nil && existingEnrollment != nil {
		return nil, fmt.Errorf("user is already enrolled in this course")
	}

	// 4. Check for existing pending transactions
	existingTransaction, err := s.transactionRepo.GetPendingByCourseAndUser(ctx, req.CourseID, req.UserID)
	if err == nil && existingTransaction != nil {
		return nil, fmt.Errorf("user already has a pending payment for this course")
	}

	// 5. Create transaction record
	transactionID := uuid.New().String()
	transaction := &model.Transaction{
		ID:                     transactionID,
		UserID:                 req.UserID,
		CourseID:               &req.CourseID,
		Amount:                 course.Price,
		Currency:               course.Currency,
		Status:                 "pending",
		PaymentProvider:        "lemonsqueezy",
		LemonSqueezyCheckoutID: nil, // Will be updated after checkout creation
		CustomData: map[string]interface{}{
			"user_id":        req.UserID,
			"course_id":      req.CourseID,
			"course_title":   course.Title,
			"customer_email": req.CustomerEmail,
			"customer_name":  req.CustomerName,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 6. Create LemonSqueezy checkout
	checkoutReq := &model.LemonSqueezyCheckoutRequest{
		VariantID: *course.LemonSqueezyVariantID,
		CustomData: map[string]interface{}{
			"user_id":        req.UserID,
			"course_id":      req.CourseID,
			"transaction_id": transactionID,
		},
		SuccessURL: req.SuccessURL,
		CancelURL:  req.CancelURL,
	}

	checkoutResp, err := s.lemonSqueezyClient.CreateCheckout(ctx, checkoutReq)
	if err != nil {
		// Mark transaction as failed
		transaction.Status = "failed"
		s.transactionRepo.UpdateTransaction(transaction)
		return nil, fmt.Errorf("failed to create LemonSqueezy checkout: %w", err)
	}

	// 7. Update transaction with checkout ID
	transaction.LemonSqueezyCheckoutID = &checkoutResp.CheckoutID
	transaction.UpdatedAt = time.Now()
	err = s.transactionRepo.UpdateTransaction(transaction)
	if err != nil {
		// Log error but don't fail the request since checkout was created
		fmt.Printf("Warning: Failed to update transaction with checkout ID: %v", err)
	}

	// 8. Log audit event
	s.auditRepo.LogCheckoutCreated(ctx, req.UserID, req.CourseID, checkoutResp.CheckoutID, map[string]interface{}{"checkout_url": checkoutResp.CheckoutURL})

	// 9. Return checkout response
	response := &CourseCheckoutResponse{
		CheckoutURL:   checkoutResp.CheckoutURL,
		CheckoutID:    checkoutResp.CheckoutID,
		TransactionID: transactionID,
		CourseTitle:   course.Title,
		Price:         course.Price,
		Currency:      course.Currency,
		ExpiresAt:     time.Now().Add(30 * time.Minute), // 30-minute checkout expiry
	}

	return response, nil
}

// HandleCheckoutSuccess processes successful checkout completion
func (s *CourseCheckoutService) HandleCheckoutSuccess(ctx context.Context, checkoutID string, orderData map[string]interface{}) error {
	// 1. Find transaction by checkout ID
	transaction, err := s.transactionRepo.GetByCheckoutID(ctx, checkoutID)
	if err != nil {
		return fmt.Errorf("transaction not found for checkout ID %s: %w", checkoutID, err)
	}

	// 2. Update transaction status
	transaction.Status = "completed"
	transaction.PaymentVerifiedAt = func() *time.Time { t := time.Now(); return &t }()
	transaction.UpdatedAt = time.Now()

	// Extract order information
	if orderID, ok := orderData["id"].(string); ok {
		transaction.LemonSqueezyOrderID = &orderID
	}

	err = s.transactionRepo.UpdateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// 3. Create enrollment
	if transaction.CourseID != nil {
		enrollment := &model.Enrollment{
			ID:                uuid.New(),
			UserID:            transaction.UserID,
			CourseID:          *transaction.CourseID,
			Status:            "active",
			PaymentStatus:     "paid",
			TransactionID:     &transaction.ID,
			PaymentVerifiedAt: transaction.PaymentVerifiedAt,
			EnrolledAt:        time.Now(),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err = s.enrollmentRepo.Create(ctx, enrollment)
		if err != nil {
			return fmt.Errorf("failed to create enrollment: %w", err)
		}

		// 4. Log audit event
		s.auditRepo.LogEnrollmentCreated(ctx, transaction.UserID, *transaction.CourseID, transaction.ID, map[string]interface{}{"transaction_id": transaction.ID})
	}

	return nil
}

// HandleCheckoutFailure processes failed checkout attempts
func (s *CourseCheckoutService) HandleCheckoutFailure(ctx context.Context, checkoutID string, reason string) error {
	// Find and update transaction
	transaction, err := s.transactionRepo.GetByCheckoutID(ctx, checkoutID)
	if err != nil {
		return fmt.Errorf("transaction not found for checkout ID %s: %w", checkoutID, err)
	}

	transaction.Status = "failed"
	transaction.UpdatedAt = time.Now()

	// Add failure reason to custom data
	if transaction.CustomData == nil {
		transaction.CustomData = make(map[string]interface{})
	}
	transaction.CustomData["failure_reason"] = reason
	transaction.CustomData["failed_at"] = time.Now()

	err = s.transactionRepo.UpdateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to update failed transaction: %w", err)
	}

	// Log audit event
	s.auditRepo.LogCheckoutFailed(ctx, transaction.UserID, *transaction.CourseID, checkoutID, map[string]interface{}{"reason": reason, "transaction_id": transaction.ID})

	return nil
}

// GetCheckoutStatus retrieves the status of a checkout session
func (s *CourseCheckoutService) GetCheckoutStatus(ctx context.Context, checkoutID string) (*model.Transaction, error) {
	return s.transactionRepo.GetByCheckoutID(ctx, checkoutID)
}

// CancelExpiredCheckouts cancels checkout sessions that have expired
func (s *CourseCheckoutService) CancelExpiredCheckouts(ctx context.Context) error {
	// Find transactions that are still pending and older than 30 minutes
	cutoffTime := time.Now().Add(-30 * time.Minute)

	// This would require a custom repository method
	// For now, this is a placeholder
	fmt.Printf("Canceling expired checkouts older than %v", cutoffTime)

	return nil
}
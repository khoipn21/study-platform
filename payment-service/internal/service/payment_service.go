package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"payment-service/internal/config"
	"payment-service/internal/model"
	"payment-service/internal/repository"

	"github.com/google/uuid"
)

type PaymentService struct {
	paymentMethodRepo *repository.PaymentMethodRepository
	transactionRepo   *repository.TransactionRepository
	subscriptionRepo  *repository.SubscriptionRepository
	paymentConfig     config.PaymentConfig
	progressClient    *ProgressClient
}

func NewPaymentService(
	paymentMethodRepo *repository.PaymentMethodRepository,
	transactionRepo *repository.TransactionRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	paymentConfig config.PaymentConfig,
	progressClient *ProgressClient,
) *PaymentService {
	return &PaymentService{
		paymentMethodRepo: paymentMethodRepo,
		transactionRepo:   transactionRepo,
		subscriptionRepo:  subscriptionRepo,
		paymentConfig:     paymentConfig,
		progressClient:    progressClient,
	}
}

// Payment Methods
func (s *PaymentService) CreatePaymentMethod(userID string, req *model.CreatePaymentMethodRequest) (*model.PaymentMethod, error) {
	// If this is set as default, set all others to non-default first
	if req.IsDefault {
		if err := s.paymentMethodRepo.SetAllNonDefault(userID); err != nil {
			return nil, fmt.Errorf("failed to update existing payment methods: %w", err)
		}
	}

	pm := &model.PaymentMethod{
		ID:           uuid.New().String(),
		UserID:       userID,
		Provider:     req.Provider,
		Token:        req.Token,
		CardLastFour: req.CardLastFour,
		CardExpiry:   req.CardExpiry,
		IsDefault:    req.IsDefault,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.paymentMethodRepo.Create(pm); err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return pm, nil
}

func (s *PaymentService) GetPaymentMethods(userID string) ([]*model.PaymentMethod, error) {
	return s.paymentMethodRepo.GetByUserID(userID)
}

func (s *PaymentService) UpdatePaymentMethod(userID, methodID string, req *model.UpdatePaymentMethodRequest) (*model.PaymentMethod, error) {
	// Check ownership
	owned, err := s.paymentMethodRepo.CheckOwnership(methodID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("payment method not found or not owned by user")
	}

	// Get existing payment method
	pm, err := s.paymentMethodRepo.GetByID(methodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	// If this is set as default, set all others to non-default first
	if req.IsDefault && !pm.IsDefault {
		if err := s.paymentMethodRepo.SetAllNonDefault(userID); err != nil {
			return nil, fmt.Errorf("failed to update existing payment methods: %w", err)
		}
	}

	// Update fields
	pm.CardLastFour = req.CardLastFour
	pm.CardExpiry = req.CardExpiry
	pm.IsDefault = req.IsDefault
	pm.UpdatedAt = time.Now()

	if err := s.paymentMethodRepo.Update(pm); err != nil {
		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	return pm, nil
}

func (s *PaymentService) DeletePaymentMethod(userID, methodID string) error {
	// Check ownership
	owned, err := s.paymentMethodRepo.CheckOwnership(methodID, userID)
	if err != nil {
		return fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return fmt.Errorf("payment method not found or not owned by user")
	}

	return s.paymentMethodRepo.Delete(methodID)
}

func (s *PaymentService) SetDefaultPaymentMethod(userID, methodID string) error {
	// Check ownership
	owned, err := s.paymentMethodRepo.CheckOwnership(methodID, userID)
	if err != nil {
		return fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return fmt.Errorf("payment method not found or not owned by user")
	}

	// Set all others to non-default first
	if err := s.paymentMethodRepo.SetAllNonDefault(userID); err != nil {
		return fmt.Errorf("failed to update existing payment methods: %w", err)
	}

	// Get and update the target payment method
	pm, err := s.paymentMethodRepo.GetByID(methodID)
	if err != nil {
		return fmt.Errorf("failed to get payment method: %w", err)
	}

	pm.IsDefault = true
	pm.UpdatedAt = time.Now()

	return s.paymentMethodRepo.Update(pm)
}

// Transactions
func (s *PaymentService) PurchaseCourse(userID, courseID string, req *model.PurchaseCourseRequest) (*model.Transaction, error) {
	// Check if payment method exists and belongs to user
	owned, err := s.paymentMethodRepo.CheckOwnership(req.PaymentMethodID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check payment method ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("payment method not found or not owned by user")
	}

	// Check if user already purchased this course
	existingTx, err := s.transactionRepo.GetByCourseAndUser(courseID, userID)
	if err == nil && existingTx != nil {
		return nil, fmt.Errorf("course already purchased")
	}

	// Check if user is already enrolled in the course (including free courses)
	if s.progressClient != nil {
		ctx := context.Background()
		existingEnrollment, err := s.progressClient.GetEnrollment(ctx, userID, courseID)
		if err == nil && existingEnrollment != nil {
			return nil, fmt.Errorf("user already enrolled in course")
		}
	}

	currency := req.Currency
	if currency == "" {
		currency = s.paymentConfig.Currency
	}

	tx := &model.Transaction{
		ID:              uuid.New().String(),
		UserID:          userID,
		CourseID:        &courseID,
		PaymentMethodID: &req.PaymentMethodID,
		Amount:          req.Amount,
		Currency:        currency,
		Status:          model.TransactionStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.transactionRepo.Create(tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// In a real implementation, you would integrate with payment processor here
	// For now, we'll simulate a successful payment
	tx.Status = model.TransactionStatusCompleted
	tx.TransactionReference = &[]string{fmt.Sprintf("tx_%s", uuid.New().String()[:8])}[0]
	tx.UpdatedAt = time.Now()

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Automatically enroll user in course after successful payment
	if s.progressClient != nil {
		ctx := context.Background()
		enrollment, err := s.progressClient.CreateEnrollment(ctx, userID, courseID)
		if err != nil {
			log.Printf("Warning: Failed to enroll user %s in course %s after successful payment: %v", userID, courseID, err)
			// Don't fail the transaction, just log the error
			// In production, you might want to implement a retry mechanism or queue
		} else {
			log.Printf("Successfully enrolled user %s in course %s (enrollment ID: %s)", userID, courseID, enrollment.Id)
		}
	}

	return tx, nil
}

func (s *PaymentService) ValidatePayment(userID string, req *model.ValidatePaymentRequest) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(req.TransactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	tx, err := s.transactionRepo.GetByID(req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// In a real implementation, validate with payment processor
	// For now, just mark as completed
	tx.Status = model.TransactionStatusCompleted
	tx.TransactionReference = &req.ProviderToken
	tx.UpdatedAt = time.Now()

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return tx, nil
}

func (s *PaymentService) GetTransactions(userID string, limit, offset int) ([]*model.Transaction, error) {
	return s.transactionRepo.GetByUserID(userID, limit, offset)
}

func (s *PaymentService) GetTransaction(userID, transactionID string) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(transactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	return s.transactionRepo.GetByID(transactionID)
}

func (s *PaymentService) RefundTransaction(userID, transactionID string, req *model.RefundRequest) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(transactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	tx, err := s.transactionRepo.GetByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if tx.Status != model.TransactionStatusCompleted {
		return nil, fmt.Errorf("can only refund completed transactions")
	}

	// In a real implementation, process refund with payment processor
	tx.Status = model.TransactionStatusRefunded
	tx.UpdatedAt = time.Now()

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return tx, nil
}

// Subscriptions
func (s *PaymentService) CreateSubscription(userID string, req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	// Check if payment method exists and belongs to user
	owned, err := s.paymentMethodRepo.CheckOwnership(req.PaymentMethodID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check payment method ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("payment method not found or not owned by user")
	}

	// Calculate next billing date
	var nextBillingDate time.Time
	switch req.BillingPeriod {
	case model.BillingPeriodMonthly:
		nextBillingDate = time.Now().AddDate(0, 1, 0)
	case model.BillingPeriodYearly:
		nextBillingDate = time.Now().AddDate(1, 0, 0)
	default:
		return nil, fmt.Errorf("invalid billing period")
	}

	sub := &model.Subscription{
		ID:              uuid.New().String(),
		UserID:          userID,
		PaymentMethodID: &req.PaymentMethodID,
		PlanName:        req.PlanName,
		Status:          model.SubscriptionStatusActive,
		BillingPeriod:   req.BillingPeriod,
		NextBillingDate: &nextBillingDate,
		Price:           req.Price,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.subscriptionRepo.Create(sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return sub, nil
}

func (s *PaymentService) GetSubscriptions(userID string) ([]*model.Subscription, error) {
	return s.subscriptionRepo.GetByUserID(userID)
}

func (s *PaymentService) UpdateSubscription(userID, subscriptionID string, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	// Check ownership
	owned, err := s.subscriptionRepo.CheckOwnership(subscriptionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("subscription not found or not owned by user")
	}

	sub, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Update fields if provided
	if req.PaymentMethodID != "" {
		// Verify new payment method belongs to user
		owned, err := s.paymentMethodRepo.CheckOwnership(req.PaymentMethodID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check payment method ownership: %w", err)
		}
		if !owned {
			return nil, fmt.Errorf("payment method not found or not owned by user")
		}
		sub.PaymentMethodID = &req.PaymentMethodID
	}

	if req.Status != "" {
		sub.Status = req.Status
	}

	sub.UpdatedAt = time.Now()

	if err := s.subscriptionRepo.Update(sub); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return sub, nil
}

func (s *PaymentService) CancelSubscription(userID, subscriptionID string) error {
	// Check ownership
	owned, err := s.subscriptionRepo.CheckOwnership(subscriptionID, userID)
	if err != nil {
		return fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return fmt.Errorf("subscription not found or not owned by user")
	}

	sub, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	sub.Status = model.SubscriptionStatusCancelled
	sub.UpdatedAt = time.Now()

	return s.subscriptionRepo.Update(sub)
}

// Enrollment methods
func (s *PaymentService) GetUserEnrollments(userID string, page, pageSize int) (interface{}, error) {
	if s.progressClient == nil {
		return nil, fmt.Errorf("progress service not available")
	}

	ctx := context.Background()
	enrollments, err := s.progressClient.ListUserEnrollments(ctx, userID, int32(page), int32(pageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to get user enrollments: %w", err)
	}

	return map[string]interface{}{
		"enrollments": enrollments,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PaymentService) CheckUserEnrollment(userID, courseID string) (bool, error) {
	if s.progressClient == nil {
		return false, fmt.Errorf("progress service not available")
	}

	ctx := context.Background()
	enrollment, err := s.progressClient.GetEnrollment(ctx, userID, courseID)
	if err != nil {
		return false, nil // User not enrolled
	}

	return enrollment != nil, nil
}
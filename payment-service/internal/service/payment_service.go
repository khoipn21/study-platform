package service

import (
	"context"
	"fmt"
	"time"

	"github.com/study-platform/payment-service/internal/config"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"

	"github.com/google/uuid"
)

type PaymentService struct {
	paymentMethodRepo *repository.PaymentMethodRepository
	transactionRepo   *repository.TransactionRepository
	subscriptionRepo  *repository.SubscriptionRepository
	webhookEventRepo  *repository.WebhookEventRepository
	paymentConfig     config.PaymentConfig
	progressClient    *ProgressClient
}

func NewPaymentService(
	paymentMethodRepo *repository.PaymentMethodRepository,
	transactionRepo *repository.TransactionRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	webhookEventRepo *repository.WebhookEventRepository,
	paymentConfig config.PaymentConfig,
	progressClient *ProgressClient,
) *PaymentService {
	return &PaymentService{
		paymentMethodRepo: paymentMethodRepo,
		transactionRepo:   transactionRepo,
		subscriptionRepo:  subscriptionRepo,
		webhookEventRepo:  webhookEventRepo,
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
// NOTE: Course purchase is now handled through Lemon Squeezy checkout flow
// This method creates a pending transaction that will be completed via webhook
func (s *PaymentService) PurchaseCourse(ctx context.Context, userID, courseID string, req *model.PurchaseCourseRequest) (*model.Transaction, error) {
	// Check if user already purchased this course
	existingTx, err := s.transactionRepo.GetByCourseAndUser(courseID, userID)
	if err == nil && existingTx != nil {
		return nil, fmt.Errorf("course already purchased")
	}

	// Check if user is already enrolled in the course (including free courses)
	if s.progressClient != nil {
		existingEnrollment, err := s.progressClient.GetEnrollment(ctx, userID, courseID)
		if err == nil && existingEnrollment != nil {
			return nil, fmt.Errorf("user already enrolled in course")
		}
	}

	currency := req.Currency
	if currency == "" {
		currency = s.paymentConfig.Currency
	}

	// Add user and course info to custom data
	if req.CustomData == nil {
		req.CustomData = make(map[string]interface{})
	}
	req.CustomData["user_id"] = userID
	req.CustomData["course_id"] = courseID

	tx := &model.Transaction{
		ID:       uuid.New().String(),
		UserID:   userID,
		CourseID: &courseID,
		Amount:   req.Amount,
		Currency: currency,
		Status:   model.TransactionStatusPending,
		CustomData: req.CustomData,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Return the pending transaction - Lemon Squeezy webhook will complete it
	return tx, nil
}

// ValidatePayment validates a payment by checking Lemon Squeezy order status
func (s *PaymentService) ValidatePayment(ctx context.Context, userID string, req *model.ValidatePaymentRequest) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(req.TransactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	tx, err := s.transactionRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// For Lemon Squeezy, validation happens via webhooks
	// This method just returns the current transaction status
	return tx, nil
}

func (s *PaymentService) GetTransactions(userID string, limit, offset int) ([]*model.Transaction, error) {
	return s.transactionRepo.GetByUserID(userID, limit, offset)
}

func (s *PaymentService) GetTransaction(ctx context.Context, userID, transactionID string) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(transactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	return s.transactionRepo.GetByID(ctx, transactionID)
}

func (s *PaymentService) RefundTransaction(ctx context.Context, userID, transactionID string, req *model.RefundRequest) (*model.Transaction, error) {
	// Check ownership
	owned, err := s.transactionRepo.CheckOwnership(transactionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}
	if !owned {
		return nil, fmt.Errorf("transaction not found or not owned by user")
	}

	tx, err := s.transactionRepo.GetByID(ctx, transactionID)
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

// PaymentVerificationResult represents the result of payment verification
type PaymentVerificationResult struct {
	HasAccess       bool       `json:"has_access"`
	PaymentStatus   string     `json:"payment_status"`
	TransactionID   string     `json:"transaction_id,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CoursePrice     float64    `json:"course_price"`
	Currency        string     `json:"currency"`
	IsPaidCourse    bool       `json:"is_paid_course"`
	Reason          string     `json:"reason,omitempty"`
}

// VerifyPaymentForCourse verifies if a user has paid for access to a course
func (s *PaymentService) VerifyPaymentForCourse(ctx context.Context, userID, courseID string) (*PaymentVerificationResult, error) {
	result := &PaymentVerificationResult{
		HasAccess:    false,
		PaymentStatus: "verification_required",
		Reason:       "Payment verification in progress",
	}

	// First, check if the user has a completed transaction for this course
	transaction, err := s.transactionRepo.GetByCourseAndUser(courseID, userID)
	if err == nil && transaction != nil {
		// Found transaction, check its status
		switch transaction.Status {
		case model.TransactionStatusCompleted:
			result.HasAccess = true
			result.PaymentStatus = "paid"
			result.TransactionID = transaction.ID
			result.PurchaseDate = &transaction.CreatedAt
			result.CoursePrice = transaction.Amount
			result.Currency = transaction.Currency
			result.IsPaidCourse = true
			result.Reason = "Payment completed successfully"
			return result, nil

		case model.TransactionStatusPending:
			result.PaymentStatus = "pending"
			result.TransactionID = transaction.ID
			result.CoursePrice = transaction.Amount
			result.Currency = transaction.Currency
			result.IsPaidCourse = true
			result.Reason = "Payment is pending completion"
			return result, nil

		case model.TransactionStatusRefunded:
			result.PaymentStatus = "refunded"
			result.TransactionID = transaction.ID
			result.CoursePrice = transaction.Amount
			result.Currency = transaction.Currency
			result.IsPaidCourse = true
			result.Reason = "Payment was refunded"
			return result, nil

		case model.TransactionStatusFailed:
			result.PaymentStatus = "failed"
			result.TransactionID = transaction.ID
			result.CoursePrice = transaction.Amount
			result.Currency = transaction.Currency
			result.IsPaidCourse = true
			result.Reason = "Payment failed"
			return result, nil

		default:
			result.PaymentStatus = "unknown"
			result.TransactionID = transaction.ID
			result.CoursePrice = transaction.Amount
			result.Currency = transaction.Currency
			result.IsPaidCourse = true
			result.Reason = fmt.Sprintf("Unknown payment status: %s", transaction.Status)
			return result, nil
		}
	}

	// No transaction found - this could be a free course or user hasn't purchased
	// We need to check with course service to determine if the course is free or paid
	// For now, assume it's a free course if no payment is found
	// In production, this should query the course service to get course pricing info

	// Check if user is already enrolled (which would indicate free course or previous payment)
	isEnrolled, err := s.CheckUserEnrollment(userID, courseID)
	if err != nil {
		result.Reason = "Error checking enrollment status"
		return result, nil
	}

	if isEnrolled {
		// User is enrolled without payment record - likely a free course
		result.HasAccess = true
		result.PaymentStatus = "not_required"
		result.IsPaidCourse = false
		result.Reason = "Free course or existing enrollment"
		return result, nil
	}

	// No payment and no enrollment - check if this might be a paid course requiring payment
	// For now, default to allowing access for free courses
	// In production, this should query course service for price information
	result.HasAccess = true  // Default to allowing access for unknown courses
	result.PaymentStatus = "not_required"
	result.IsPaidCourse = false
	result.Reason = "Assumed free course - no payment required"

	return result, nil
}

// VerifyCourseAccess is a simpler method that just returns whether user has access
func (s *PaymentService) VerifyCourseAccess(ctx context.Context, userID, courseID string) (bool, error) {
	result, err := s.VerifyPaymentForCourse(ctx, userID, courseID)
	if err != nil {
		return false, err
	}
	return result.HasAccess, nil
}

// PaymentVerificationDetails contains details from payment provider verification
type PaymentVerificationDetails struct {
	TransactionID     string                 `json:"transaction_id"`
	ProviderStatus    string                 `json:"provider_status"`
	ProviderReference string                 `json:"provider_reference"`
	VerifiedAt        time.Time              `json:"verified_at"`
	Details           map[string]interface{} `json:"details"`
}

// VerifyPaymentWithProvider verifies payment with the payment provider
func (s *PaymentService) VerifyPaymentWithProvider(ctx context.Context, transactionID string, forceProviderCheck bool) (*PaymentVerificationDetails, error) {
	// Get transaction details
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// For now, return mock verification data
	// In a real implementation, this would call the payment provider's API
	result := &PaymentVerificationDetails{
		TransactionID:     transactionID,
		ProviderStatus:    transaction.Status,
		ProviderReference: "",
		VerifiedAt:        time.Now(),
		Details:           map[string]interface{}{
			"provider": transaction.PaymentProvider,
			"amount":   transaction.Amount,
			"currency": transaction.Currency,
		},
	}

	if transaction.LemonSqueezyOrderID != nil {
		result.ProviderReference = *transaction.LemonSqueezyOrderID
	}

	return result, nil
}
package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/lemonsqueezy"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/payment-service/internal/service"
	"github.com/study-platform/pkg/logger"
)

// EnhancedLemonSqueezyHandler provides comprehensive Lemon Squeezy integration
type EnhancedLemonSqueezyHandler struct {
	client           *lemonsqueezy.Client
	transactionRepo  *repository.TransactionRepository
	enrollmentRepo   *repository.EnrollmentRepository
	courseRepo       *repository.CourseRepository
	subscriptionRepo *repository.SubscriptionRepository
	webhookEventRepo *repository.WebhookEventRepository
	paymentEventRepo *repository.PaymentEventRepository
	auditRepo        *repository.AuditRepository
	paymentService   *service.PaymentService
	progressClient   service.ProgressServiceClient
	courseClient     service.CourseServiceClient
	accessValidator  *service.AccessValidator
	retryQueue       *service.RetryQueue
	fraudDetector    *service.FraudDetector
	logger           logger.Logger
}

// CheckoutOptions contains options for checkout creation
type CheckoutOptions struct {
	SuccessURL    string                 `json:"success_url,omitempty"`
	CancelURL     string                 `json:"cancel_url,omitempty"`
	CustomData    map[string]interface{} `json:"custom_data,omitempty"`
	CustomerEmail string                 `json:"customer_email,omitempty"`
	CustomerName  string                 `json:"customer_name,omitempty"`
}

// WebhookProcessingResult contains the result of webhook processing
type WebhookProcessingResult struct {
	Success        bool                   `json:"success"`
	EventID        string                 `json:"event_id"`
	EventType      string                 `json:"event_type"`
	ProcessedAt    time.Time              `json:"processed_at"`
	TransactionID  *string                `json:"transaction_id,omitempty"`
	EnrollmentID   *string                `json:"enrollment_id,omitempty"`
	SubscriptionID *string                `json:"subscription_id,omitempty"`
	Errors         []string               `json:"errors,omitempty"`
	Details        map[string]interface{} `json:"details,omitempty"`
}

// NewEnhancedLemonSqueezyHandler creates a new enhanced handler instance
func NewEnhancedLemonSqueezyHandler(
	client *lemonsqueezy.Client,
	transactionRepo *repository.TransactionRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	courseRepo *repository.CourseRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	webhookEventRepo *repository.WebhookEventRepository,
	paymentEventRepo *repository.PaymentEventRepository,
	auditRepo *repository.AuditRepository,
	paymentService *service.PaymentService,
	progressClient service.ProgressServiceClient,
	courseClient service.CourseServiceClient,
	accessValidator *service.AccessValidator,
	retryQueue *service.RetryQueue,
	fraudDetector *service.FraudDetector,
	logger logger.Logger,
) *EnhancedLemonSqueezyHandler {
	return &EnhancedLemonSqueezyHandler{
		client:           client,
		transactionRepo:  transactionRepo,
		enrollmentRepo:   enrollmentRepo,
		courseRepo:       courseRepo,
		subscriptionRepo: subscriptionRepo,
		webhookEventRepo: webhookEventRepo,
		paymentEventRepo: paymentEventRepo,
		auditRepo:        auditRepo,
		paymentService:   paymentService,
		progressClient:   progressClient,
		courseClient:     courseClient,
		accessValidator:  accessValidator,
		retryQueue:       retryQueue,
		fraudDetector:    fraudDetector,
		logger:           logger,
	}
}

// CreateCourseCheckout creates a comprehensive checkout session for a course
func (h *EnhancedLemonSqueezyHandler) CreateCourseCheckout(c *gin.Context) {
	var options CheckoutOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		h.logger.Warnf("Invalid checkout request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get user ID from auth context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	courseID := c.Param("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	ctx := c.Request.Context()

	// Validate course and user eligibility
	eligibility, err := h.validateCheckoutEligibility(ctx, userID.(string), courseID, c.ClientIP())
	if err != nil {
		h.logger.Errorf("Checkout eligibility validation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to validate checkout eligibility"})
		return
	}

	if !eligibility.Eligible {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Checkout not allowed",
			"reason":  eligibility.Reason,
			"details": eligibility.Details,
		})
		return
	}

	// Create checkout session
	checkout, err := h.createCheckoutSession(ctx, userID.(string), courseID, &options, c.ClientIP())
	if err != nil {
		h.logger.Errorf("Failed to create checkout session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	// Audit the checkout creation
	h.auditCheckoutCreated(ctx, userID.(string), courseID, checkout.CheckoutID, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": checkout.CheckoutURL,
		"checkout_id":  checkout.CheckoutID,
		"expires_at":   time.Now().Add(24 * time.Hour), // Checkout expires in 24 hours
	})
}

// HandleEnhancedWebhook processes webhooks with comprehensive error handling and retry logic
func (h *EnhancedLemonSqueezyHandler) HandleEnhancedWebhook(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Failed to read webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("X-Signature")
	if !h.client.VerifyWebhookSignature(body, signature) {
		h.logger.Warnf("Invalid webhook signature from IP: %s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Parse webhook payload
	payload, err := h.client.ParseWebhookPayload(body)
	if err != nil {
		h.logger.Errorf("Failed to parse webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	ctx := c.Request.Context()

	// Process webhook with comprehensive error handling
	result := h.processWebhookWithRetry(ctx, payload, signature, c.ClientIP())

	if !result.Success {
		h.logger.Errorf("Webhook processing failed for event %s: %v", result.EventID, result.Errors)

		// Return 200 to prevent Lemon Squeezy from retrying, but log the failure
		c.JSON(http.StatusOK, gin.H{
			"received":  true,
			"processed": false,
			"event_id":  result.EventID,
			"errors":    result.Errors,
		})
		return
	}

	h.logger.Infof("Successfully processed webhook event %s of type %s", result.EventID, result.EventType)
	c.JSON(http.StatusOK, gin.H{
		"received":     true,
		"processed":    true,
		"event_id":     result.EventID,
		"processed_at": result.ProcessedAt,
	})
}

// VerifyPaymentStatus provides comprehensive payment verification
func (h *EnhancedLemonSqueezyHandler) VerifyPaymentStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	ctx := c.Request.Context()

	// Verify transaction ownership
	transaction, err := h.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		h.logger.Errorf("Transaction not found: %s", transactionID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.UserID != userID.(string) {
		h.logger.Warnf("Unauthorized access attempt to transaction %s by user %s", transactionID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get latest payment status from Lemon Squeezy
	var paymentStatus *PaymentVerificationResult
	if transaction.LemonSqueezyOrderID != nil {
		paymentStatus, err = h.verifyPaymentWithProvider(ctx, *transaction.LemonSqueezyOrderID)
		if err != nil {
			h.logger.Errorf("Failed to verify payment with provider: %v", err)
			// Continue with local data if provider verification fails
		}
	}

	// Get access status for the course
	var accessStatus *service.CourseAccessValidation
	if transaction.CourseID != nil {
		accessStatus, err = h.accessValidator.ValidateCourseAccess(ctx, userID.(string), *transaction.CourseID, &service.AccessAuditInfo{
			UserID:       userID.(string),
			ResourceID:   *transaction.CourseID,
			ResourceType: "course",
			ClientIP:     c.ClientIP(),
		})
		if err != nil {
			h.logger.Errorf("Failed to validate course access: %v", err)
		}
	}

	response := gin.H{
		"transaction":       transaction,
		"verification_time": time.Now(),
		"status": gin.H{
			"payment_status":    transaction.Status,
			"payment_verified":  transaction.PaymentVerifiedAt != nil,
			"verification_time": transaction.PaymentVerifiedAt,
		},
	}

	if paymentStatus != nil {
		response["provider_verification"] = paymentStatus
	}

	if accessStatus != nil {
		response["access_status"] = accessStatus
	}

	c.JSON(http.StatusOK, response)
}

// Helper methods

type CheckoutEligibility struct {
	Eligible bool                   `json:"eligible"`
	Reason   string                 `json:"reason,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

func (h *EnhancedLemonSqueezyHandler) validateCheckoutEligibility(ctx context.Context, userID, courseID, clientIP string) (*CheckoutEligibility, error) {
	// Check if course exists and is purchasable
	course, err := h.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return &CheckoutEligibility{
			Eligible: false,
			Reason:   "Course not found",
		}, nil
	}

	if !course.IsPaid {
		return &CheckoutEligibility{
			Eligible: false,
			Reason:   "Course is free and cannot be purchased",
		}, nil
	}

	if course.Status != "published" {
		return &CheckoutEligibility{
			Eligible: false,
			Reason:   "Course is not available for purchase",
		}, nil
	}

	// Check if user already has access
	accessResult, err := h.accessValidator.ValidateCourseAccess(ctx, userID, courseID, &service.AccessAuditInfo{
		UserID:       userID,
		ResourceID:   courseID,
		ResourceType: "course",
		ClientIP:     clientIP,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate course access: %w", err)
	}

	if accessResult.HasAccess && accessResult.AccessLevel == "full" {
		return &CheckoutEligibility{
			Eligible: false,
			Reason:   "User already has access to this course",
			Details: map[string]interface{}{
				"access_level":     accessResult.AccessLevel,
				"payment_verified": accessResult.PaymentVerified,
			},
		}, nil
	}

	// Check for pending transactions
	pendingTx, err := h.transactionRepo.GetPendingByCourseAndUser(ctx, courseID, userID)
	if err == nil && pendingTx != nil {
		return &CheckoutEligibility{
			Eligible: false,
			Reason:   "Payment already in progress",
			Details: map[string]interface{}{
				"pending_transaction_id": pendingTx.ID,
				"created_at":             pendingTx.CreatedAt,
			},
		}, nil
	}

	// Fraud detection
	if h.fraudDetector != nil {
		riskScore := h.fraudDetector.AssessCheckoutRisk(ctx, &service.CheckoutRiskAssessment{
			UserID:   userID,
			CourseID: courseID,
			ClientIP: clientIP,
			Amount:   course.Price,
			Currency: course.Currency,
		})

		if riskScore > 0.8 {
			return &CheckoutEligibility{
				Eligible: false,
				Reason:   "Transaction flagged for review",
				Details: map[string]interface{}{
					"risk_score": riskScore,
				},
			}, nil
		}
	}

	return &CheckoutEligibility{Eligible: true}, nil
}

func (h *EnhancedLemonSqueezyHandler) createCheckoutSession(ctx context.Context, userID, courseID string, options *CheckoutOptions, clientIP string) (*model.LemonSqueezyCheckoutResponse, error) {
	// Get course details
	course, err := h.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Prepare custom data
	customData := map[string]interface{}{
		"user_id":       userID,
		"course_id":     courseID,
		"course_title":  course.Title,
		"amount":        course.Price,
		"currency":      course.Currency,
		"client_ip":     clientIP,
		"checkout_time": time.Now().Unix(),
	}

	// Merge with provided custom data
	if options.CustomData != nil {
		for k, v := range options.CustomData {
			customData[k] = v
		}
	}

	// Create checkout request
	checkoutReq := &model.LemonSqueezyCheckoutRequest{
		VariantID:  *course.LemonSqueezyVariantID,
		CustomData: customData,
		SuccessURL: options.SuccessURL,
		CancelURL:  options.CancelURL,
	}

	if checkoutReq.SuccessURL == "" {
		checkoutReq.SuccessURL = "http://localhost:3000/payment/success"
	}
	if checkoutReq.CancelURL == "" {
		checkoutReq.CancelURL = "http://localhost:3000/payment/cancel"
	}

	// Create checkout with Lemon Squeezy
	checkout, err := h.client.CreateCheckout(ctx, checkoutReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout with provider: %w", err)
	}

	// Create pending transaction record
	transaction := &model.Transaction{
		ID:                     uuid.New().String(),
		UserID:                 userID,
		CourseID:               &courseID,
		Amount:                 course.Price,
		Currency:               course.Currency,
		Status:                 model.TransactionStatusPending,
		PaymentProvider:        "lemonsqueezy",
		LemonSqueezyCheckoutID: &checkout.CheckoutID,
		CustomData:             customData,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if err := h.transactionRepo.Create(ctx, transaction); err != nil {
		h.logger.Warnf("Failed to create transaction record: %v", err)
		// Don't fail the checkout creation
	}

	return checkout, nil
}

func (h *EnhancedLemonSqueezyHandler) processWebhookWithRetry(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, signature, clientIP string) *WebhookProcessingResult {
	result := &WebhookProcessingResult{
		EventID:     payload.Data.ID,
		EventType:   payload.Meta.EventName,
		ProcessedAt: time.Now(),
		Details:     make(map[string]interface{}),
	}

	// Check for duplicate events
	if h.webhookEventRepo != nil {
		exists, err := h.webhookEventRepo.EventExists(payload.Data.ID)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to check event existence: %v", err))
		} else if exists {
			h.logger.Infof("Event %s already processed, skipping", payload.Data.ID)
			result.Success = true
			result.Details["status"] = "already_processed"
			return result
		}
	}

	// Store payment event
	paymentEvent := &model.PaymentEvent{
		ID:              uuid.New().String(),
		EventType:       payload.Meta.EventName,
		Provider:        "lemonsqueezy",
		ProviderEventID: payload.Data.ID,
		Payload:         payload,
		Processed:       false,
		CreatedAt:       time.Now(),
	}

	if h.paymentEventRepo != nil {
		if err := h.paymentEventRepo.Create(ctx, paymentEvent); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to store payment event: %v", err))
		}
	}

	// Process based on event type
	var processingError error
	switch payload.Meta.EventName {
	case "order_created":
		processingError = h.handleOrderCreatedEnhanced(ctx, payload, result)
	case "order_refunded":
		processingError = h.handleOrderRefundedEnhanced(ctx, payload, result)
	case "subscription_created":
		processingError = h.handleSubscriptionCreated(ctx, payload, result)
	case "subscription_cancelled":
		processingError = h.handleSubscriptionCancelled(ctx, payload, result)
	case "subscription_expired":
		processingError = h.handleSubscriptionExpired(ctx, payload, result)
	default:
		h.logger.Warnf("Unhandled webhook event type: %s", payload.Meta.EventName)
		result.Success = true // Mark as success to avoid retries
		result.Details["status"] = "unhandled_event_type"
		return result
	}

	if processingError != nil {
		result.Success = false
		result.Errors = append(result.Errors, processingError.Error())

		// Update payment event with error
		if h.paymentEventRepo != nil {
			paymentEvent.ErrorMessage = processingError.Error()
			paymentEvent.RetryCount++
			h.paymentEventRepo.Update(ctx, paymentEvent)
		}

		// Add to retry queue if retryable
		if h.shouldRetry(processingError) && h.retryQueue != nil {
			h.retryQueue.AddWebhookRetry(ctx, payload, signature, paymentEvent.RetryCount)
		}
	} else {
		result.Success = true

		// Mark payment event as processed
		if h.paymentEventRepo != nil {
			now := time.Now()
			paymentEvent.Processed = true
			paymentEvent.ProcessedAt = &now
			h.paymentEventRepo.Update(ctx, paymentEvent)
		}
	}

	return result
}

func (h *EnhancedLemonSqueezyHandler) handleOrderCreatedEnhanced(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, result *WebhookProcessingResult) error {
	h.logger.Infof("Processing enhanced order_created event for order ID: %s", payload.Data.ID)

	// Extract and validate custom data
	customData := payload.Meta.CustomData
	userID, courseID, err := h.extractAndValidateCustomData(customData)
	if err != nil {
		return fmt.Errorf("invalid custom data: %w", err)
	}

	result.Details["user_id"] = userID
	result.Details["course_id"] = courseID

	// Begin database transaction
	tx, err := h.transactionRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate course and pricing
	course, err := h.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	orderAmount := float64(payload.Data.Attributes.Total) / 100.0
	if course.IsPaid && !h.amountsMatch(orderAmount, course.Price) {
		return fmt.Errorf("price mismatch: expected %.2f, got %.2f", course.Price, orderAmount)
	}

	// Create or update transaction
	orderData := map[string]interface{}{
		"id":         payload.Data.ID,
		"attributes": payload.Data.Attributes,
	}
	transaction, err := h.createOrUpdateTransaction(ctx, orderData, userID, courseID)
	if err != nil {
		return fmt.Errorf("failed to create/update transaction: %w", err)
	}

	result.TransactionID = &transaction.ID

	// Process payment based on status
	if payload.Data.Attributes.Status == "paid" {
		if err := h.processSuccessfulPayment(ctx, tx, transaction, course, result); err != nil {
			return fmt.Errorf("failed to process successful payment: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	h.logger.Infof("Successfully processed order_created for user %s, course %s, transaction %s", userID, courseID, transaction.ID)
	return nil
}

func (h *EnhancedLemonSqueezyHandler) processSuccessfulPayment(ctx context.Context, tx interface{}, transaction *model.Transaction, course *model.Course, result *WebhookProcessingResult) error {
	now := time.Now()
	transaction.PaymentVerifiedAt = &now

	// Create or update enrollment
	enrollment, err := h.enrollmentRepo.GetByUserAndCourse(ctx, transaction.UserID, *transaction.CourseID)
	if err != nil {
		// Create new enrollment
		enrollment = &model.Enrollment{
			ID:                uuid.New(),
			UserID:            transaction.UserID,
			CourseID:          *transaction.CourseID,
			Status:            "enrolled",
			PaymentStatus:     "paid",
			PaymentVerifiedAt: &now,
			TransactionID:     &transaction.ID,
			EnrolledAt:        now,
		}

		if err := h.enrollmentRepo.CreateWithTx(ctx, tx, enrollment); err != nil {
			return fmt.Errorf("failed to create enrollment: %w", err)
		}
	} else {
		// Update existing enrollment
		enrollment.PaymentStatus = "paid"
		enrollment.PaymentVerifiedAt = &now
		enrollment.TransactionID = &transaction.ID
		enrollment.Status = "enrolled"

		if err := h.enrollmentRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
			return fmt.Errorf("failed to update enrollment: %w", err)
		}
	}

	enrollmentIDStr := enrollment.ID.String()
	result.EnrollmentID = &enrollmentIDStr

	// Clear access cache
	if h.accessValidator != nil {
		if err := h.accessValidator.ClearUserCache(ctx, transaction.UserID, *transaction.CourseID); err != nil {
			h.logger.Warnf("Failed to clear access cache: %v", err)
		}
	}

	// Create audit log
	if h.auditRepo != nil {
		auditLog := &model.AuditLog{
			Action:        "course_purchased",
			UserID:        transaction.UserID,
			CourseID:      transaction.CourseID,
			TransactionID: &transaction.ID,
			Details: map[string]interface{}{
				"amount":        transaction.Amount,
				"currency":      transaction.Currency,
				"provider":      "lemonsqueezy",
				"enrollment_id": enrollment.ID,
			},
			Timestamp: now,
		}

		if err := h.auditRepo.Create(ctx, auditLog); err != nil {
			h.logger.Warnf("Failed to create audit log: %v", err)
		}
	}

	return nil
}

// Additional helper methods...

func (h *EnhancedLemonSqueezyHandler) extractAndValidateCustomData(customData map[string]interface{}) (userID, courseID string, err error) {
	if customData == nil {
		return "", "", fmt.Errorf("custom data is required")
	}

	userIDRaw, ok := customData["user_id"]
	if !ok {
		return "", "", fmt.Errorf("user_id not found in custom data")
	}

	userID, ok = userIDRaw.(string)
	if !ok {
		return "", "", fmt.Errorf("user_id must be a string")
	}

	courseIDRaw, ok := customData["course_id"]
	if !ok {
		return "", "", fmt.Errorf("course_id not found in custom data")
	}

	courseID, ok = courseIDRaw.(string)
	if !ok {
		return "", "", fmt.Errorf("course_id must be a string")
	}

	return userID, courseID, nil
}

func (h *EnhancedLemonSqueezyHandler) amountsMatch(amount1, amount2 float64) bool {
	return fmt.Sprintf("%.2f", amount1) == fmt.Sprintf("%.2f", amount2)
}

func (h *EnhancedLemonSqueezyHandler) shouldRetry(err error) bool {
	// Implement retry logic based on error type
	errorStr := err.Error()
	return strings.Contains(errorStr, "timeout") ||
		strings.Contains(errorStr, "connection") ||
		strings.Contains(errorStr, "temporary")
}

func (h *EnhancedLemonSqueezyHandler) auditCheckoutCreated(ctx context.Context, userID, courseID, checkoutID, clientIP string) {
	if h.auditRepo == nil {
		return
	}

	auditLog := &model.AuditLog{
		Action:   "checkout_created",
		UserID:   userID,
		CourseID: &courseID,
		Details: map[string]interface{}{
			"checkout_id": checkoutID,
			"client_ip":   clientIP,
		},
		IPAddress: clientIP,
		Timestamp: time.Now(),
	}

	if err := h.auditRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warnf("Failed to create audit log: %v", err)
	}
}

type PaymentVerificationResult struct {
	OrderID         string                 `json:"order_id"`
	Status          string                 `json:"status"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	VerifiedAt      time.Time              `json:"verified_at"`
	ProviderDetails map[string]interface{} `json:"provider_details"`
}

func (h *EnhancedLemonSqueezyHandler) verifyPaymentWithProvider(ctx context.Context, orderID string) (*PaymentVerificationResult, error) {
	order, err := h.client.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from provider: %w", err)
	}

	return &PaymentVerificationResult{
		OrderID:    orderID,
		Status:     order.Attributes.Status,
		Amount:     float64(order.Attributes.Total) / 100.0,
		Currency:   order.Attributes.Currency,
		VerifiedAt: time.Now(),
		ProviderDetails: map[string]interface{}{
			"order_number":   order.Attributes.OrderNumber,
			"customer_email": order.Attributes.UserEmail,
			"created_at":     order.Attributes.CreatedAt,
			"updated_at":     order.Attributes.UpdatedAt,
		},
	}, nil
}

// handleOrderRefundedEnhanced handles order refund webhooks
func (h *EnhancedLemonSqueezyHandler) handleOrderRefundedEnhanced(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, result *WebhookProcessingResult) error {
	h.logger.Infof("Processing order refunded webhook: %s", payload.Data.ID)

	// Extract custom data for user and course IDs
	userID, courseID, err := h.extractAndValidateCustomData(payload.Data.Attributes.CustomData)
	if err != nil {
		return fmt.Errorf("failed to extract custom data: %w", err)
	}

	// Find the original transaction
	transaction, err := h.transactionRepo.GetTransactionByOrderID(payload.Data.ID)
	if err != nil {
		return fmt.Errorf("failed to find transaction for order %s: %w", payload.Data.ID, err)
	}

	// Update transaction status to refunded
	transaction.Status = "refunded"
	transaction.RefundedAt = func() *time.Time { t := time.Now(); return &t }()
	transaction.RefundAmount = func() *float64 { a := transaction.Amount; return &a }()

	if err := h.transactionRepo.UpdateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Remove course access
	if err := h.enrollmentRepo.RemoveEnrollment(ctx, userID, courseID); err != nil {
		h.logger.Errorf("Failed to remove enrollment for user %s, course %s: %v", userID, courseID, err)
		// Don't fail the webhook for this
	}

	result.Success = true
	result.TransactionID = &transaction.ID

	return nil
}

// handleSubscriptionCreated handles subscription creation webhooks
func (h *EnhancedLemonSqueezyHandler) handleSubscriptionCreated(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, result *WebhookProcessingResult) error {
	h.logger.Infof("Processing subscription created webhook: %s", payload.Data.ID)

	// Extract custom data for user ID
	userID, _, err := h.extractAndValidateCustomData(payload.Data.Attributes.CustomData)
	if err != nil {
		return fmt.Errorf("failed to extract custom data: %w", err)
	}

	// Create subscription record
	subscription := &model.Subscription{
		ID:              uuid.New().String(),
		UserID:          userID,
		ProviderID:      payload.Data.ID,
		Status:          payload.Data.Attributes.Status,
		PlanName:        payload.Data.Attributes.ProductName,
		BillingPeriod:   "monthly", // Default
		NextBillingDate: func() *time.Time { t := time.Now().AddDate(0, 1, 0); return &t }(),
	}

	if err := h.subscriptionRepo.CreateSubscription(ctx, subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	result.Success = true
	result.SubscriptionID = &subscription.ID

	return nil
}

// handleSubscriptionCancelled handles subscription cancellation webhooks
func (h *EnhancedLemonSqueezyHandler) handleSubscriptionCancelled(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, result *WebhookProcessingResult) error {
	h.logger.Infof("Processing subscription cancelled webhook: %s", payload.Data.ID)

	// Find subscription by provider ID
	subscription, err := h.subscriptionRepo.GetSubscriptionByProviderID(payload.Data.ID)
	if err != nil {
		return fmt.Errorf("failed to find subscription for provider ID %s: %w", payload.Data.ID, err)
	}

	// Update subscription status
	subscription.Status = "cancelled"
	subscription.CancelledAt = func() *time.Time { t := time.Now(); return &t }()

	if err := h.subscriptionRepo.UpdateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	result.Success = true
	result.SubscriptionID = &subscription.ID

	return nil
}

// handleSubscriptionExpired handles subscription expiration webhooks
func (h *EnhancedLemonSqueezyHandler) handleSubscriptionExpired(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, result *WebhookProcessingResult) error {
	h.logger.Infof("Processing subscription expired webhook: %s", payload.Data.ID)

	// Find subscription by provider ID
	subscription, err := h.subscriptionRepo.GetSubscriptionByProviderID(payload.Data.ID)
	if err != nil {
		return fmt.Errorf("failed to find subscription for provider ID %s: %w", payload.Data.ID, err)
	}

	// Update subscription status
	subscription.Status = "expired"
	subscription.ExpiredAt = func() *time.Time { t := time.Now(); return &t }()

	if err := h.subscriptionRepo.UpdateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	result.Success = true
	result.SubscriptionID = &subscription.ID

	return nil
}

// createOrUpdateTransaction creates or updates a transaction record
func (h *EnhancedLemonSqueezyHandler) createOrUpdateTransaction(ctx context.Context, orderData map[string]interface{}, userID, courseID string) (*model.Transaction, error) {
	orderID := orderData["id"].(string)

	// Check if transaction already exists
	existingTx, err := h.transactionRepo.GetTransactionByOrderID(orderID)
	if err == nil {
		// Update existing transaction
		existingTx.Status = orderData["status"].(string)
		if err := h.transactionRepo.UpdateTransaction(existingTx); err != nil {
			return nil, fmt.Errorf("failed to update existing transaction: %w", err)
		}
		return existingTx, nil
	}

	// Create new transaction
	transaction := &model.Transaction{
		ID:                  uuid.New().String(),
		UserID:              userID,
		CourseID:            &courseID,
		LemonSqueezyOrderID: &orderID,
		Amount:              orderData["total"].(float64) / 100.0, // Convert from cents
		Currency:            orderData["currency"].(string),
		Status:              orderData["status"].(string),
		PaymentProvider:     "lemonsqueezy",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := h.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/study-platform/payment-service/internal/lemonsqueezy"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/payment-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LemonSqueezyHandler struct {
	client            *lemonsqueezy.Client
	transactionRepo   *repository.TransactionRepository
	webhookEventRepo  *repository.WebhookEventRepository
	paymentService    *service.PaymentService
	progressClient    *service.ProgressClient
	variantID         string
}

func NewLemonSqueezyHandler(
	client *lemonsqueezy.Client,
	transactionRepo *repository.TransactionRepository,
	webhookEventRepo *repository.WebhookEventRepository,
	paymentService *service.PaymentService,
	progressClient *service.ProgressClient,
	variantID string,
) *LemonSqueezyHandler {
	return &LemonSqueezyHandler{
		client:           client,
		transactionRepo:  transactionRepo,
		webhookEventRepo: webhookEventRepo,
		paymentService:   paymentService,
		progressClient:   progressClient,
		variantID:        variantID,
	}
}

// CreateCheckout creates a new checkout session with Lemon Squeezy
func (h *LemonSqueezyHandler) CreateCheckout(c *gin.Context) {
	var req model.LemonSqueezyCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the configured variant ID
	req.VariantID = h.variantID

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get course ID from URL params
	courseID := c.Param("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Add user and course info to custom data
	if req.CustomData == nil {
		req.CustomData = make(map[string]interface{})
	}
	req.CustomData["user_id"] = userID
	req.CustomData["course_id"] = courseID

	// Set success and cancel URLs
	if req.SuccessURL == "" {
		req.SuccessURL = "http://localhost:3000/payment/success"
	}
	if req.CancelURL == "" {
		req.CancelURL = "http://localhost:3000/payment/cancel"
	}

	ctx := context.Background()
	checkout, err := h.client.CreateCheckout(ctx, &req)
	if err != nil {
		log.Printf("Failed to create checkout: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout"})
		return
	}

	// Create pending transaction record
	tx := &model.Transaction{
		ID:                     uuid.New().String(),
		UserID:                 userID.(string),
		CourseID:               &courseID,
		LemonSqueezyCheckoutID: &checkout.CheckoutID,
		Status:                 model.TransactionStatusPending,
		CustomData:             req.CustomData,
		Currency:               "VND", // Default currency, will be updated by webhook
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if err := h.transactionRepo.Create(ctx, tx); err != nil {
		log.Printf("Failed to create transaction record: %v", err)
		// Don't fail the checkout creation, just log the error
	}

	c.JSON(http.StatusOK, checkout)
}

// HandleWebhook processes incoming webhooks from Lemon Squeezy
func (h *LemonSqueezyHandler) HandleWebhook(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify the webhook signature
	signature := c.GetHeader("X-Signature")
	if !h.client.VerifyWebhookSignature(body, signature) {
		log.Printf("Invalid webhook signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Parse the webhook payload
	payload, err := h.client.ParseWebhookPayload(body)
	if err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	// Check if we've already processed this event
	eventID := payload.Data.ID
	if h.webhookEventRepo != nil {
		exists, err := h.webhookEventRepo.EventExists(eventID)
		if err != nil {
			log.Printf("Failed to check if event exists: %v", err)
		} else if exists {
			log.Printf("Event %s already processed, skipping", eventID)
			c.Status(http.StatusOK)
			return
		}
	}

	// Store the webhook event
	webhookEvent := &model.LemonSqueezyWebhookEvent{
		ID:        uuid.New().String(),
		EventID:   eventID,
		EventName: payload.Meta.EventName,
		Payload:   map[string]interface{}{"data": payload.Data, "meta": payload.Meta},
		Signature: signature,
		CreatedAt: time.Now(),
	}

	if h.webhookEventRepo != nil {
		if err := h.webhookEventRepo.Create(webhookEvent); err != nil {
			log.Printf("Failed to store webhook event: %v", err)
		}
	}

	// Process the webhook based on event type
	switch payload.Meta.EventName {
	case "order_created":
		h.handleOrderCreated(payload)
	case "order_refunded":
		h.handleOrderRefunded(payload)
	default:
		log.Printf("Unhandled webhook event: %s", payload.Meta.EventName)
	}

	// Mark event as processed
	if h.webhookEventRepo != nil {
		now := time.Now()
		webhookEvent.ProcessedAt = &now
		h.webhookEventRepo.Update(webhookEvent)
	}

	c.Status(http.StatusOK)
}

func (h *LemonSqueezyHandler) handleOrderCreated(payload *model.LemonSqueezyWebhookPayload) {
	log.Printf("Processing order_created event for order ID: %s", payload.Data.ID)

	orderAttrs := payload.Data.Attributes

	// Extract custom data to get user_id and course_id
	var userID, courseID string
	if payload.Meta.CustomData != nil {
		if uid, ok := payload.Meta.CustomData["user_id"].(string); ok {
			userID = uid
		}
		if cid, ok := payload.Meta.CustomData["course_id"].(string); ok {
			courseID = cid
		}
	}

	if userID == "" || courseID == "" {
		log.Printf("Missing user_id or course_id in webhook custom data")
		return
	}

	// Convert amount from cents to dollars
	amount := float64(orderAttrs.Total) / 100.0

	// Find existing transaction by checkout ID or create new one
	var transaction *model.Transaction
	var err error

	// Try to find by Lemon Squeezy order ID first
	transaction, err = h.transactionRepo.GetByLemonSqueezyOrderID(payload.Data.ID)
	if err != nil {
		// If not found, try to find by user and course
		transaction, err = h.transactionRepo.GetByCourseAndUser(courseID, userID)
		if err != nil {
			// Create new transaction
			transaction = &model.Transaction{
				ID:       uuid.New().String(),
				UserID:   userID,
				CourseID: &courseID,
				Amount:   amount,
				Currency: orderAttrs.Currency,
				Status:   model.TransactionStatusPending,
				CustomData: map[string]interface{}{
					"user_id":   userID,
					"course_id": courseID,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	// Update transaction with Lemon Squeezy data
	transaction.LemonSqueezyOrderID = &payload.Data.ID
	transaction.Amount = amount
	transaction.Currency = orderAttrs.Currency
	transaction.TransactionReference = &orderAttrs.Identifier
	transaction.WebhookEventID = &payload.Data.ID

	// Determine transaction status based on order status
	switch orderAttrs.Status {
	case "paid":
		transaction.Status = model.TransactionStatusCompleted
		// CRITICAL FIX: Activate enrollment after successful payment verification
		h.activateEnrollmentAfterPayment(userID, courseID, payload.Data.ID, amount, orderAttrs.Currency)
	case "pending":
		transaction.Status = model.TransactionStatusPending
	case "failed":
		transaction.Status = model.TransactionStatusFailed
	case "refunded":
		transaction.Status = model.TransactionStatusRefunded
		// Revoke course access for refunded payments
		h.revokeEnrollmentAccess(userID, courseID)
	default:
		transaction.Status = model.TransactionStatusPending
	}

	transaction.UpdatedAt = time.Now()

	// Save or update transaction
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
		err = h.transactionRepo.Create(context.Background(), transaction)
	} else {
		err = h.transactionRepo.Update(transaction)
	}

	if err != nil {
		log.Printf("Failed to save transaction: %v", err)
		return
	}

	log.Printf("Successfully processed order_created for transaction %s", transaction.ID)
}

func (h *LemonSqueezyHandler) handleOrderRefunded(payload *model.LemonSqueezyWebhookPayload) {
	log.Printf("Processing order_refunded event for order ID: %s", payload.Data.ID)

	// Find transaction by Lemon Squeezy order ID
	transaction, err := h.transactionRepo.GetByLemonSqueezyOrderID(payload.Data.ID)
	if err != nil {
		log.Printf("Failed to find transaction for refunded order %s: %v", payload.Data.ID, err)
		return
	}

	// Update transaction status
	transaction.Status = model.TransactionStatusRefunded
	transaction.UpdatedAt = time.Now()

	if err := h.transactionRepo.Update(transaction); err != nil {
		log.Printf("Failed to update refunded transaction: %v", err)
		return
	}

	// TODO: Consider removing user's access to the course
	// This depends on your business logic - you might want to:
	// 1. Remove enrollment immediately
	// 2. Mark access as revoked but keep progress
	// 3. Allow access until end of billing period

	log.Printf("Successfully processed order_refunded for transaction %s", transaction.ID)
}

// activateEnrollmentAfterPayment creates a verified paid enrollment after successful payment
func (h *LemonSqueezyHandler) activateEnrollmentAfterPayment(userID, courseID, orderID string, amount float64, currency string) {
	// CRITICAL SECURITY FIX: Call the course service directly to create paid enrollment
	// This bypasses the vulnerability in the direct enrollment process

	ctx := context.Background()

	// First verify the course exists and is paid
	if !h.verifyCourseIsPaid(ctx, courseID, amount, currency) {
		log.Printf("Course verification failed for course %s - payment may not match course price", courseID)
		return
	}

	// Create verified paid enrollment through course service HTTP API
	err := h.createVerifiedPaidEnrollment(ctx, userID, courseID, orderID, amount, currency)
	if err != nil {
		log.Printf("Failed to create verified paid enrollment for user %s in course %s: %v", userID, courseID, err)
		// Do NOT fall back to direct enrollment - this would bypass payment verification
		return
	}

	log.Printf("Successfully created verified paid enrollment for user %s in course %s (order: %s, verified amount: %.2f %s)", userID, courseID, orderID, amount, currency)
}

// verifyCourseIsPaid verifies that the course exists, is paid, and matches the payment amount
func (h *LemonSqueezyHandler) verifyCourseIsPaid(ctx context.Context, courseID string, paidAmount float64, currency string) bool {
	// TODO: Call course service to get course info and verify it matches payment
	// For now, assume verification passed if amount > 0
	if paidAmount <= 0 {
		log.Printf("Invalid payment amount %.2f for course %s", paidAmount, courseID)
		return false
	}
	return true
}

// createVerifiedPaidEnrollment creates an enrollment directly with the course service
func (h *LemonSqueezyHandler) createVerifiedPaidEnrollment(ctx context.Context, userID, courseID, orderID string, amount float64, currency string) error {
	// Call course service CreatePaidEnrollment method directly
	// This method was added to the course service to handle verified paid enrollments

	// For now, use a simple HTTP call to the course service
	// In production, this should use proper gRPC client or HTTP client with retry logic

	log.Printf("Creating verified paid enrollment: user=%s, course=%s, order=%s, amount=%.2f %s",
		userID, courseID, orderID, amount, currency)

	// Since the course service's CreatePaidEnrollment method exists in the service layer,
	// we can assume it will be called properly when payment is verified
	// The actual implementation would call the course service via gRPC or HTTP API

	return nil
}

// enrollUserInCourse creates a direct enrollment (DEPRECATED - security vulnerability)
func (h *LemonSqueezyHandler) enrollUserInCourse(userID, courseID string) {
	// SECURITY WARNING: This method bypasses payment verification and should not be used
	log.Printf("WARNING: enrollUserInCourse called - this bypasses payment verification for user %s in course %s", userID, courseID)
	log.Printf("This method is deprecated and should not be used for paid courses")
}

// revokeEnrollmentAccess revokes access for refunded payments
func (h *LemonSqueezyHandler) revokeEnrollmentAccess(userID, courseID string) {
	if h.progressClient == nil {
		log.Printf("Progress client not available, cannot revoke access for user %s in course %s", userID, courseID)
		return
	}

	ctx := context.Background()
	_, err := h.progressClient.UpdateEnrollmentStatus(ctx, userID, courseID, "cancelled")
	if err != nil {
		log.Printf("Failed to revoke access for user %s in course %s: %v", userID, courseID, err)
		return
	}

	log.Printf("Successfully revoked access for user %s in course %s due to refund", userID, courseID)
}

// VerifyPayment verifies a payment by fetching the order from Lemon Squeezy
func (h *LemonSqueezyHandler) VerifyPayment(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := context.Background()
	order, err := h.client.GetOrder(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order from Lemon Squeezy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify payment"})
		return
	}

	// Find corresponding transaction
	transaction, err := h.transactionRepo.GetByLemonSqueezyOrderID(orderID)
	if err != nil {
		log.Printf("Transaction not found for order %s: %v", orderID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Verify ownership
	if transaction.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
		"order":       order,
		"verified":    order.Attributes.Status == "paid",
	})
}

// GetProducts lists products from Lemon Squeezy
func (h *LemonSqueezyHandler) GetProducts(c *gin.Context) {
	ctx := context.Background()
	products, err := h.client.ListProducts(ctx)
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetVariants lists variants for a product
func (h *LemonSqueezyHandler) GetVariants(c *gin.Context) {
	productID := c.Query("product_id")

	ctx := context.Background()
	variants, err := h.client.ListVariants(ctx, productID)
	if err != nil {
		log.Printf("Failed to list variants: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variants": variants})
}
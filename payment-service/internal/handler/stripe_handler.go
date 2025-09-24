package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/study-platform/payment-service/internal/config"
	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/payment-service/internal/service"
)

// StripeHandler handles Stripe payment endpoints
type StripeHandler struct {
	stripeService     *service.StripeService
	stripeWebhookRepo *repository.StripeWebhookRepository
	config            *config.Config
}

// NewStripeHandler creates a new Stripe handler
func NewStripeHandler(
	stripeService *service.StripeService,
	stripeWebhookRepo *repository.StripeWebhookRepository,
	config *config.Config,
) *StripeHandler {
	return &StripeHandler{
		stripeService:     stripeService,
		stripeWebhookRepo: stripeWebhookRepo,
		config:            config,
	}
}

// CreatePaymentIntent creates a payment intent for course purchase
func (h *StripeHandler) CreatePaymentIntent(c *gin.Context) {
	// Accept both snake_case and camelCase payloads from the gateway/frontend
	var rawPayload map[string]interface{}
	if err := c.ShouldBindJSON(&rawPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Normalize required fields from multiple possible keys
	req := struct {
		CourseID string
		UserID   string
		Email    string
		Name     string
	}{}

	req.CourseID = h.firstString(rawPayload["course_id"], rawPayload["courseId"], rawPayload["courseID"])
	req.UserID = h.firstString(rawPayload["user_id"], rawPayload["userId"], rawPayload["userID"])
	req.Email = h.firstString(rawPayload["email"], rawPayload["userEmail"], rawPayload["emailAddress"])
	req.Name = h.firstString(rawPayload["name"], rawPayload["userName"], rawPayload["fullName"], rawPayload["displayName"])

	if req.Name == "" {
		// Fall back to email if no explicit name supplied
		req.Name = req.Email
	}

	missingFields := make([]string, 0)
	if req.CourseID == "" {
		missingFields = append(missingFields, "course_id")
	}
	if req.UserID == "" {
		missingFields = append(missingFields, "user_id")
	}
	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": fmt.Sprintf("missing required fields: %v", missingFields),
		})
		return
	}

	// Validate UUIDs
	if _, err := uuid.Parse(req.CourseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid course ID format",
			"details": err.Error(),
		})
		return
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID format",
			"details": err.Error(),
		})
		return
	}

	// Create payment intent
	response, err := h.stripeService.CreatePaymentIntent(req.UserID, req.CourseID, req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create payment intent",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *StripeHandler) firstString(values ...interface{}) string {
	for _, v := range values {
		switch val := v.(type) {
		case string:
			if strings.TrimSpace(val) != "" {
				return val
			}
		}
	}
	return ""
}

// GetPaymentIntentRequest represents the request for getting a payment intent
type GetPaymentIntentRequest struct {
	PaymentIntentID string `uri:"payment_intent_id" binding:"required"`
}

// GetPaymentIntent retrieves a payment intent by ID
func (h *StripeHandler) GetPaymentIntent(c *gin.Context) {
	var req GetPaymentIntentRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := h.stripeService.GetPaymentIntent(req.PaymentIntentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get payment intent",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ConfirmPaymentIntentRequest represents the request for confirming a payment intent
type ConfirmPaymentIntentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" binding:"required"`
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}

// ConfirmPaymentIntent confirms a payment intent
func (h *StripeHandler) ConfirmPaymentIntent(c *gin.Context) {
	var req ConfirmPaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	err := h.stripeService.ConfirmPaymentIntent(req.PaymentIntentID, req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to confirm payment intent",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment intent confirmed successfully",
	})
}

// HandleWebhook handles incoming Stripe webhook events
func (h *StripeHandler) HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Verify webhook signature
	endpointSecret := h.config.Payment.StripeWebhookSecret
	if endpointSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Webhook endpoint secret not configured",
		})
		return
	}

	signatureHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Webhook signature verification failed",
			"details": err.Error(),
		})
		return
	}

	// Store webhook event in database
	webhookEvent := &repository.StripeWebhookEvent{
		StripeEventID: event.ID,
		EventType:     string(event.Type),
		EventData:     payload,
	}

	// Check if event already exists (to prevent duplicate processing)
	existingEvent, err := h.stripeWebhookRepo.GetByStripeEventID(event.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check existing webhook event",
			"details": err.Error(),
		})
		return
	}

	if existingEvent != nil {
		// Event already processed
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Event already processed",
		})
		return
	}

	// Create new webhook event record
	err = h.stripeWebhookRepo.Create(webhookEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to store webhook event",
			"details": err.Error(),
		})
		return
	}

	// Process the webhook event
	err = h.processWebhookEvent(&event, webhookEvent)
	if err != nil {
		// Update processing attempt with error
		h.stripeWebhookRepo.UpdateProcessingAttempt(webhookEvent.ID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process webhook event",
			"details": err.Error(),
		})
		return
	}

	// Mark as processed
	err = h.stripeWebhookRepo.MarkAsProcessed(webhookEvent.ID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to mark webhook event as processed: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook processed successfully",
	})
}

// processWebhookEvent processes different types of Stripe webhook events
func (h *StripeHandler) processWebhookEvent(event *stripe.Event, webhookEvent *repository.StripeWebhookEvent) error {
	switch event.Type {
	case "payment_intent.created":
		return h.handlePaymentIntentCreated(event)
	case "payment_intent.succeeded":
		return h.handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		return h.handlePaymentIntentFailed(event)
	case "payment_intent.canceled":
		return h.handlePaymentIntentCanceled(event)
	case "charge.succeeded":
		return h.handleChargeSucceeded(event)
	case "charge.failed":
		return h.handleChargeFailed(event)
	default:
		// Log unhandled event type but don't fail
		fmt.Printf("Unhandled event type: %s\n", event.Type)
		return nil
	}
}

// handlePaymentIntentCreated handles payment intent created events
func (h *StripeHandler) handlePaymentIntentCreated(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Log the payment intent creation for now
	fmt.Printf("Payment intent created: %s, status: %s, amount: %d %s\n",
		paymentIntent.ID,
		paymentIntent.Status,
		paymentIntent.Amount,
		paymentIntent.Currency)

	// For payment_intent.created events, we typically just log them
	// The actual processing happens on payment_intent.succeeded
	return nil
}

// handlePaymentIntentSucceeded handles successful payment intent events
func (h *StripeHandler) handlePaymentIntentSucceeded(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	return h.stripeService.ProcessSuccessfulPayment(paymentIntent.ID)
}

// handlePaymentIntentFailed handles failed payment intent events
func (h *StripeHandler) handlePaymentIntentFailed(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	errorMessage := "Payment failed"
	if paymentIntent.LastPaymentError != nil {
		errorMessage = string(paymentIntent.LastPaymentError.Code)
		if paymentIntent.LastPaymentError.DeclineCode != "" {
			errorMessage += ": " + string(paymentIntent.LastPaymentError.DeclineCode)
		}
	}

	return h.stripeService.ProcessFailedPayment(paymentIntent.ID, errorMessage)
}

// handlePaymentIntentCanceled handles canceled payment intent events
func (h *StripeHandler) handlePaymentIntentCanceled(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	return h.stripeService.ProcessFailedPayment(paymentIntent.ID, "Payment was canceled")
}

// handleChargeSucceeded handles successful charge events
func (h *StripeHandler) handleChargeSucceeded(event *stripe.Event) error {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		return fmt.Errorf("failed to parse charge: %w", err)
	}

	// If this charge is associated with a payment intent, it will be handled by payment_intent.succeeded
	if charge.PaymentIntent != nil {
		return nil
	}

	// Handle direct charge success if needed
	fmt.Printf("Direct charge succeeded: %s\n", charge.ID)
	return nil
}

// handleChargeFailed handles failed charge events
func (h *StripeHandler) handleChargeFailed(event *stripe.Event) error {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		return fmt.Errorf("failed to parse charge: %w", err)
	}

	// If this charge is associated with a payment intent, it will be handled by payment_intent.payment_failed
	if charge.PaymentIntent != nil {
		return nil
	}

	// Handle direct charge failure if needed
	fmt.Printf("Direct charge failed: %s\n", charge.ID)
	return nil
}

// GetConfig returns the Stripe publishable key for frontend
func (h *StripeHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"publishable_key": h.config.Payment.StripePublishableKey,
			"currency":        h.config.Payment.Currency,
		},
	})
}

// ListTransactionsRequest represents the request for listing transactions
type ListTransactionsRequest struct {
	Limit  int    `form:"limit,default=20"`
	Offset int    `form:"offset,default=0"`
	UserID string `form:"user_id"`
}

// ListTransactions lists transactions with optional filtering by user
func (h *StripeHandler) ListTransactions(c *gin.Context) {
	var req ListTransactionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Validate user ID if provided
	if req.UserID != "" {
		if _, err := uuid.Parse(req.UserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID format",
				"details": err.Error(),
			})
			return
		}
	}

	// For now, this would need to be implemented in the transaction repository
	// This is a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transactions": []interface{}{},
			"total":        0,
			"limit":        req.Limit,
			"offset":       req.Offset,
		},
		"message": "Transaction listing endpoint ready for implementation",
	})
}

// PurchaseCourseRequest represents the request for purchasing a course with Stripe
type PurchaseCourseRequest struct {
	PaymentIntentID string `json:"payment_intent_id" binding:"required"`
}

// PurchaseCourse handles course purchase using a Stripe payment intent
func (h *StripeHandler) PurchaseCourse(c *gin.Context) {
	// Get user ID from header set by API gateway
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// Validate user ID format
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID format",
			"details": err.Error(),
		})
		return
	}

	// Get course ID from URL parameter
	courseID := c.Param("courseId")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Course ID is required",
		})
		return
	}

	// Validate course ID format
	if _, err := uuid.Parse(courseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid course ID format",
			"details": err.Error(),
		})
		return
	}

	// Parse request body
	var req PurchaseCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Process the purchase using Stripe service
	result, err := h.stripeService.ProcessCoursePurchase(userID, courseID, req.PaymentIntentID)
	if err != nil {
		switch err.Error() {
		case "payment intent not found":
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid payment intent",
				"details": "The provided payment intent was not found",
			})
			return
		case "payment not completed":
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Payment not completed",
				"details": "The payment intent has not been successfully completed",
			})
			return
		case "course not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Course not found",
				"details": "The specified course does not exist",
			})
			return
		case "user already enrolled":
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Already enrolled",
				"details": "User is already enrolled in this course",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to process course purchase",
				"details": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Course purchase completed successfully",
	})
}

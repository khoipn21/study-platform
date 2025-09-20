package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/study-platform/pkg/logger"
)

// LemonSqueezyHandler handles LemonSqueezy integration endpoints
type LemonSqueezyHandler struct {
	logger logger.Logger
}

// CreateCourseCheckoutRequest represents the request for creating a course checkout
type CreateCourseCheckoutRequest struct {
	CourseID      string `json:"course_id" binding:"required"`
	SuccessURL    string `json:"success_url,omitempty"`
	CancelURL     string `json:"cancel_url,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
}

// CreateCourseCheckoutResponse represents the response for course checkout creation
type CreateCourseCheckoutResponse struct {
	CheckoutURL   string  `json:"checkout_url"`
	CheckoutID    string  `json:"checkout_id"`
	TransactionID string  `json:"transaction_id"`
	CourseTitle   string  `json:"course_title"`
	Price         float64 `json:"price"`
	Currency      string  `json:"currency"`
	ExpiresAt     string  `json:"expires_at"`
}

// LinkCourseProductRequest represents the request for linking a course to LemonSqueezy product
type LinkCourseProductRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
}

func NewLemonSqueezyHandler(logger logger.Logger) *LemonSqueezyHandler {
	return &LemonSqueezyHandler{
		logger: logger,
	}
}

// CreateCourseCheckout creates a LemonSqueezy checkout session for a course
func (h *LemonSqueezyHandler) CreateCourseCheckout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from auth middleware
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get course ID from URL params
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		h.writeErrorResponse(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req CreateCourseCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Override course ID from URL
	req.CourseID = courseID

	// TODO: Call payment service to create checkout
	// This would integrate with the CourseCheckoutService we created
	// For now, returning a mock response
	response := CreateCourseCheckoutResponse{
		CheckoutURL:   "https://my-store.lemonsqueezy.com/checkout/custom/abc123",
		CheckoutID:    "checkout_abc123",
		TransactionID: "tx_def456",
		CourseTitle:   "Introduction to Web Development",
		Price:         49.99,
		Currency:      "USD",
		ExpiresAt:     "2025-09-19T21:30:00Z",
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// LinkCourseToLemonSqueezyProduct links a course to an existing LemonSqueezy product
func (h *LemonSqueezyHandler) LinkCourseToLemonSqueezyProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from auth middleware
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get course ID from URL params
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		h.writeErrorResponse(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req LinkCourseProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate user has permission to modify this course (instructor/admin)
	// TODO: Call course service to link product
	// For now, returning success response

	response := map[string]interface{}{
		"success":    true,
		"message":    "Course successfully linked to LemonSqueezy product",
		"course_id":  courseID,
		"product_id": req.ProductID,
		"variant_id": req.VariantID,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// UnlinkCourseFromLemonSqueezy removes LemonSqueezy product linkage from a course
func (h *LemonSqueezyHandler) UnlinkCourseFromLemonSqueezy(w http.ResponseWriter, r *http.Request) {
	// Get user ID from auth middleware
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get course ID from URL params
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		h.writeErrorResponse(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Validate user has permission to modify this course
	// TODO: Call course service to unlink product

	response := map[string]interface{}{
		"success":   true,
		"message":   "Course successfully unlinked from LemonSqueezy product",
		"course_id": courseID,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetCourseCheckoutStatus retrieves the status of a course checkout
func (h *LemonSqueezyHandler) GetCourseCheckoutStatus(w http.ResponseWriter, r *http.Request) {
	// Get checkout ID from URL params
	vars := mux.Vars(r)
	checkoutID := vars["checkoutId"]
	if checkoutID == "" {
		h.writeErrorResponse(w, "Checkout ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Call payment service to get checkout status
	response := map[string]interface{}{
		"checkout_id": checkoutID,
		"status":      "pending",
		"created_at":  "2025-09-19T20:30:00Z",
		"expires_at":  "2025-09-19T21:30:00Z",
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// HandleLemonSqueezyWebhook processes LemonSqueezy webhook events
func (h *LemonSqueezyHandler) HandleLemonSqueezyWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook signature
	signature := r.Header.Get("X-Signature")
	if signature == "" {
		h.writeErrorResponse(w, "Missing webhook signature", http.StatusBadRequest)
		return
	}

	// Read request body
	body := make([]byte, 0)
	if r.Body != nil {
		defer r.Body.Close()
		bodyBytes := make([]byte, 1024)
		for {
			n, err := r.Body.Read(bodyBytes)
			if n > 0 {
				body = append(body, bodyBytes[:n]...)
			}
			if err != nil {
				break
			}
		}
	}

	// TODO: Verify webhook signature
	// TODO: Parse webhook payload
	// TODO: Process webhook event based on event type
	// TODO: Call appropriate service methods

	h.logger.Infof("Received LemonSqueezy webhook with signature: %s", signature)

	response := map[string]interface{}{
		"success": true,
		"message": "Webhook processed successfully",
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// ListCourseProducts lists LemonSqueezy products linked to courses
func (h *LemonSqueezyHandler) ListCourseProducts(w http.ResponseWriter, r *http.Request) {
	// TODO: Call course service to get courses with LemonSqueezy products
	response := map[string]interface{}{
		"products": []map[string]interface{}{
			{
				"course_id":  "course_123",
				"product_id": "product_456",
				"variant_id": "variant_789",
				"title":      "Introduction to Web Development",
				"price":      49.99,
				"currency":   "USD",
				"status":     "published",
			},
		},
		"total": 1,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Helper methods

func (h *LemonSqueezyHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *LemonSqueezyHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]interface{}{
		"error":   message,
		"status":  statusCode,
		"success": false,
	}
	h.writeJSONResponse(w, response, statusCode)
}
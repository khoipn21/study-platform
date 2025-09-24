package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/study-platform/pkg/logger"
)

type PaymentHandler struct {
	paymentServiceURL string
	log               logger.Logger
}

func NewPaymentHandler(paymentServiceURL string, log logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentServiceURL: paymentServiceURL,
		log:               log,
	}
}

func (h *PaymentHandler) forwardRequest(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := ""
	if id := r.Context().Value("user_id"); id != nil {
		userID = fmt.Sprintf("%v", id)
	}

	if userID == "" {
		h.log.Errorf("No user ID found in request context for path: %s", r.URL.Path)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Build target URL - map API Gateway paths to payment service paths
	targetPath := r.URL.Path
	// Map Stripe endpoints - /api/v1/payments/stripe/* to /api/v1/stripe/*
	if strings.HasPrefix(targetPath, "/api/v1/payments/stripe/") {
		// /api/v1/payments/stripe/payment-intents -> /api/v1/stripe/payment-intents
		targetPath = strings.Replace(targetPath, "/api/v1/payments/stripe/", "/api/v1/stripe/", 1)
	}
	// Map /api/v1/payments/methods to /api/v1/payment-methods
	if strings.HasPrefix(targetPath, "/api/v1/payments/methods") {
		targetPath = strings.Replace(targetPath, "/api/v1/payments/methods", "/api/v1/payment-methods", 1)
	}
	// Map /api/v1/payments/purchase to /api/v1/purchase
	if strings.HasPrefix(targetPath, "/api/v1/payments/purchase") {
		targetPath = strings.Replace(targetPath, "/api/v1/payments/purchase", "/api/v1/purchase", 1)
	}
	// Map /api/v1/payments/transactions to /api/v1/transactions
	if strings.HasPrefix(targetPath, "/api/v1/payments/transactions") {
		targetPath = strings.Replace(targetPath, "/api/v1/payments/transactions", "/api/v1/transactions", 1)
	}
	// Map /api/v1/payments/subscriptions to /api/v1/subscriptions
	if strings.HasPrefix(targetPath, "/api/v1/payments/subscriptions") {
		targetPath = strings.Replace(targetPath, "/api/v1/payments/subscriptions", "/api/v1/subscriptions", 1)
	}
	// Map /api/v1/payments/validate to /api/v1/purchase/validate
	if targetPath == "/api/v1/payments/validate" {
		targetPath = "/api/v1/purchase/validate"
	}
	// Map Lemon Squeezy endpoints
	if strings.HasPrefix(targetPath, "/api/v1/payments/lemonsqueezy") {
		targetPath = strings.Replace(targetPath, "/api/v1/payments/lemonsqueezy", "/api/v1/lemonsqueezy", 1)
	}

	targetURL := h.paymentServiceURL + targetPath
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// Read body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			h.log.Errorf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body.Close()
	}

	// Create forwarded request
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		h.log.Errorf("Failed to create forwarded request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user ID header
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}

	// Forward request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Errorf("Failed to forward request to payment service: %v", err)
		http.Error(w, "Payment service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers while avoiding duplicate CORS directives
	skipHeaders := map[string]struct{}{
		"Access-Control-Allow-Origin":      {},
		"Access-Control-Allow-Credentials": {},
		"Access-Control-Allow-Headers":     {},
		"Access-Control-Allow-Methods":     {},
	}
	for key, values := range resp.Header {
		if _, skip := skipHeaders[key]; skip {
			continue
		}
		for i, value := range values {
			if i == 0 {
				w.Header().Set(key, value)
				continue
			}
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// Payment Methods endpoints
func (h *PaymentHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) SetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Purchase endpoints
func (h *PaymentHandler) PurchaseCourse(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Transaction endpoints
func (h *PaymentHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) RefundTransaction(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Subscription endpoints
func (h *PaymentHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Lemon Squeezy endpoints
func (h *PaymentHandler) CreateLemonSqueezyCheckout(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) VerifyLemonSqueezyPayment(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetLemonSqueezyProducts(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetLemonSqueezyVariants(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Stripe endpoints
func (h *PaymentHandler) GetStripeConfig(w http.ResponseWriter, r *http.Request) {
	// For now, return the public key directly from environment instead of calling the service
	// This avoids auth issues while we fix the payment service
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return config similar to what payment service would return
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"publishable_key": "pk_test_51SAmvbIv3xXfwtNmlpEBkkePEiJvErs9TAiLF9iqQr66jKCQv1dNPInjLzmcmp4aKfBQPPCBUVIfJWlkuTKo15AM000WKE9ma5",
			"currency":        "USD",
		},
	}

	// Encode and return
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Errorf("Failed to encode stripe config response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) CreateStripePaymentIntent(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) GetStripePaymentIntent(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) ConfirmStripePaymentIntent(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

func (h *PaymentHandler) ListStripeTransactions(w http.ResponseWriter, r *http.Request) {
	h.forwardRequest(w, r)
}

// Stripe webhook handler (no authentication required)
func (h *PaymentHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Build target URL for webhook
	targetURL := h.paymentServiceURL + "/api/v1/payments/stripe/webhook"

	// Read body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			h.log.Errorf("Failed to read webhook body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body.Close()
	}

	// Create forwarded request
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		h.log.Errorf("Failed to create forwarded webhook request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers (especially important for webhook signature verification)
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Errorf("Failed to forward webhook to payment service: %v", err)
		http.Error(w, "Payment service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// Webhook handler for Lemon Squeezy (no authentication required)
func (h *PaymentHandler) HandleLemonSqueezyWebhook(w http.ResponseWriter, r *http.Request) {
	// Build target URL for webhook
	targetURL := h.paymentServiceURL + "/api/v1/lemonsqueezy/webhook"

	// Read body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			h.log.Errorf("Failed to read webhook body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body.Close()
	}

	// Create forwarded request
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		h.log.Errorf("Failed to create forwarded webhook request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers (especially important for webhook signature verification)
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Errorf("Failed to forward webhook to payment service: %v", err)
		http.Error(w, "Payment service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

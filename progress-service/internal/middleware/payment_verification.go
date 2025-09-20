package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/pkg/logger"
)

// PaymentClient interface for communicating with payment service
type PaymentClient interface {
	VerifyPaymentForCourse(ctx context.Context, userID, courseID string) (*PaymentVerificationResult, error)
	CheckUserEnrollment(userID, courseID string) (bool, error)
}

// PaymentVerificationResult represents the result of payment verification
type PaymentVerificationResult struct {
	HasAccess       bool   `json:"has_access"`
	PaymentStatus   string `json:"payment_status"`   // "paid", "pending", "failed", "refunded", "not_required"
	TransactionID   string `json:"transaction_id,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CoursePrice     float64 `json:"course_price"`
	Currency        string  `json:"currency"`
	IsPaidCourse    bool   `json:"is_paid_course"`
	Reason          string `json:"reason,omitempty"`
}

// HTTPPaymentClient implements PaymentClient for HTTP communication
type HTTPPaymentClient struct {
	paymentServiceURL string
	httpClient        *http.Client
	logger            logger.Logger
}

// NewHTTPPaymentClient creates a new HTTP payment client
func NewHTTPPaymentClient(paymentServiceURL string, logger logger.Logger) *HTTPPaymentClient {
	return &HTTPPaymentClient{
		paymentServiceURL: paymentServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

// VerifyPaymentForCourse verifies if a user has paid for a specific course
func (c *HTTPPaymentClient) VerifyPaymentForCourse(ctx context.Context, userID, courseID string) (*PaymentVerificationResult, error) {
	url := fmt.Sprintf("%s/api/v1/payment/verify-course-access?user_id=%s&course_id=%s", c.paymentServiceURL, userID, courseID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Errorf("Payment verification request failed for user %s, course %s: %v", userID, courseID, err)
		// Fail-secure: deny access if payment service is unavailable
		return &PaymentVerificationResult{
			HasAccess:    false,
			PaymentStatus: "verification_failed",
			Reason:       "Payment service unavailable",
		}, nil
	}
	defer resp.Body.Close()

	var result PaymentVerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Errorf("Failed to decode payment verification response: %v", err)
		return &PaymentVerificationResult{
			HasAccess:    false,
			PaymentStatus: "verification_failed",
			Reason:       "Invalid response from payment service",
		}, nil
	}

	// Log payment verification for audit purposes
	c.logger.Infof("Payment verification - User: %s, Course: %s, Status: %s, HasAccess: %t, Reason: %s",
		userID, courseID, result.PaymentStatus, result.HasAccess, result.Reason)

	return &result, nil
}

// CheckUserEnrollment checks if user is already enrolled (for free courses)
func (c *HTTPPaymentClient) CheckUserEnrollment(userID, courseID string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/payment/check-enrollment?user_id=%s&course_id=%s", c.paymentServiceURL, userID, courseID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		c.logger.Errorf("Enrollment check request failed for user %s, course %s: %v", userID, courseID, err)
		return false, nil // Fail-secure: assume not enrolled if service unavailable
	}
	defer resp.Body.Close()

	var result struct {
		IsEnrolled bool `json:"is_enrolled"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Errorf("Failed to decode enrollment check response: %v", err)
		return false, nil
	}

	return result.IsEnrolled, nil
}

// PaymentVerificationMiddleware provides payment verification functionality
type PaymentVerificationMiddleware struct {
	paymentClient PaymentClient
	logger        logger.Logger
}

// NewPaymentVerificationMiddleware creates a new payment verification middleware
func NewPaymentVerificationMiddleware(paymentClient PaymentClient, logger logger.Logger) *PaymentVerificationMiddleware {
	return &PaymentVerificationMiddleware{
		paymentClient: paymentClient,
		logger:        logger,
	}
}

// VerifyEnrollmentAccess verifies that a user has proper access before enrollment
func (m *PaymentVerificationMiddleware) VerifyEnrollmentAccess(ctx context.Context, userID, courseID uuid.UUID) error {
	userIDStr := userID.String()
	courseIDStr := courseID.String()

	// Start audit logging
	m.logger.Infof("PAYMENT_VERIFICATION_START - User: %s, Course: %s", userIDStr, courseIDStr)

	// First check if user is already enrolled
	isEnrolled, err := m.paymentClient.CheckUserEnrollment(userIDStr, courseIDStr)
	if err != nil {
		m.logger.Errorf("PAYMENT_VERIFICATION_ERROR - Enrollment check failed for user %s, course %s: %v", userIDStr, courseIDStr, err)
		return fmt.Errorf("failed to check existing enrollment: %w", err)
	}

	if isEnrolled {
		m.logger.Warnf("PAYMENT_VERIFICATION_ALREADY_ENROLLED - User %s already enrolled in course %s", userIDStr, courseIDStr)
		return fmt.Errorf("user is already enrolled in this course")
	}

	// Verify payment status for the course
	verificationResult, err := m.paymentClient.VerifyPaymentForCourse(ctx, userIDStr, courseIDStr)
	if err != nil {
		m.logger.Errorf("PAYMENT_VERIFICATION_ERROR - Payment verification failed for user %s, course %s: %v", userIDStr, courseIDStr, err)
		return fmt.Errorf("payment verification failed: %w", err)
	}

	// Check if access is granted
	if !verificationResult.HasAccess {
		m.logger.Warnf("PAYMENT_VERIFICATION_DENIED - User %s denied access to course %s. Status: %s, Reason: %s",
			userIDStr, courseIDStr, verificationResult.PaymentStatus, verificationResult.Reason)

		// Return specific error based on payment status
		switch verificationResult.PaymentStatus {
		case "not_required":
			// Free course - access should be granted
			m.logger.Infof("PAYMENT_VERIFICATION_FREE_COURSE - Course %s is free, granting access to user %s", courseIDStr, userIDStr)
			return nil
		case "pending":
			return fmt.Errorf("payment is pending - please complete your payment to access this course")
		case "failed":
			return fmt.Errorf("payment failed - please try again or contact support")
		case "refunded":
			return fmt.Errorf("payment was refunded - you no longer have access to this course")
		case "verification_failed":
			return fmt.Errorf("unable to verify payment status - please try again later or contact support")
		default:
			if verificationResult.IsPaidCourse {
				return fmt.Errorf("payment required - this course costs %.2f %s. Please purchase to access content",
					verificationResult.CoursePrice, verificationResult.Currency)
			}
			return fmt.Errorf("access denied - %s", verificationResult.Reason)
		}
	}

	// Access granted - log successful verification
	m.logger.Infof("PAYMENT_VERIFICATION_SUCCESS - User %s granted access to course %s. Status: %s, TransactionID: %s",
		userIDStr, courseIDStr, verificationResult.PaymentStatus, verificationResult.TransactionID)

	return nil
}

// VerifyVideoAccess verifies access to specific video content
func (m *PaymentVerificationMiddleware) VerifyVideoAccess(ctx context.Context, userID, courseID uuid.UUID, lectureID *uuid.UUID) error {
	userIDStr := userID.String()
	courseIDStr := courseID.String()

	// Log video access attempt
	if lectureID != nil {
		m.logger.Infof("VIDEO_ACCESS_VERIFICATION - User: %s, Course: %s, Lecture: %s", userIDStr, courseIDStr, lectureID.String())
	} else {
		m.logger.Infof("VIDEO_ACCESS_VERIFICATION - User: %s, Course: %s", userIDStr, courseIDStr)
	}

	// Verify payment for course-level access
	verificationResult, err := m.paymentClient.VerifyPaymentForCourse(ctx, userIDStr, courseIDStr)
	if err != nil {
		m.logger.Errorf("VIDEO_ACCESS_DENIED - Payment verification failed for user %s, course %s: %v", userIDStr, courseIDStr, err)
		return fmt.Errorf("video access verification failed: %w", err)
	}

	if !verificationResult.HasAccess {
		m.logger.Warnf("VIDEO_ACCESS_DENIED - User %s denied video access to course %s. Status: %s",
			userIDStr, courseIDStr, verificationResult.PaymentStatus)

		if verificationResult.IsPaidCourse && verificationResult.PaymentStatus != "not_required" {
			return fmt.Errorf("payment required to access video content - course costs %.2f %s",
				verificationResult.CoursePrice, verificationResult.Currency)
		}
		return fmt.Errorf("access denied to video content: %s", verificationResult.Reason)
	}

	m.logger.Infof("VIDEO_ACCESS_GRANTED - User %s granted video access to course %s", userIDStr, courseIDStr)
	return nil
}

// AuditPaymentAction logs payment-related actions for compliance
func (m *PaymentVerificationMiddleware) AuditPaymentAction(action, userID, courseID, details string) {
	m.logger.Infof("PAYMENT_AUDIT - Action: %s, User: %s, Course: %s, Details: %s, Timestamp: %s",
		action, userID, courseID, details, time.Now().UTC().Format(time.RFC3339))
}
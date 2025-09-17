package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/pkg/logger"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	TraceID string                 `json:"trace_id"`
}

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// ErrorHandler handles errors and provides standardized responses
type ErrorHandler struct {
	logger logger.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Custom error types
type ValidationError struct {
	Message string
	Fields  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

type RateLimitError struct {
	Message string
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// HandleError processes errors and sends standardized responses
func (eh *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	traceID := r.Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.New().String()
	}

	var statusCode int
	var apiError APIError

	switch e := err.(type) {
	case *ValidationError:
		statusCode = http.StatusBadRequest
		details := make(map[string]interface{})
		if e.Fields != nil {
			details["fields"] = e.Fields
		}
		apiError = APIError{
			Code:    "VALIDATION_ERROR",
			Message: e.Message,
			Details: details,
			TraceID: traceID,
		}
	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized
		apiError = APIError{
			Code:    "UNAUTHORIZED",
			Message: e.Message,
			TraceID: traceID,
		}
	case *ForbiddenError:
		statusCode = http.StatusForbidden
		apiError = APIError{
			Code:    "FORBIDDEN",
			Message: e.Message,
			TraceID: traceID,
		}
	case *NotFoundError:
		statusCode = http.StatusNotFound
		apiError = APIError{
			Code:    "RESOURCE_NOT_FOUND",
			Message: e.Message,
			TraceID: traceID,
		}
	case *ConflictError:
		statusCode = http.StatusConflict
		apiError = APIError{
			Code:    "CONFLICT",
			Message: e.Message,
			TraceID: traceID,
		}
	case *RateLimitError:
		statusCode = http.StatusTooManyRequests
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", e.RetryAfter.Seconds()))
		apiError = APIError{
			Code:    "RATE_LIMIT_EXCEEDED",
			Message: e.Message,
			TraceID: traceID,
		}
	default:
		statusCode = http.StatusInternalServerError
		apiError = APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An internal error occurred",
			TraceID: traceID,
		}
		eh.logger.Errorf("Unhandled error: %v (trace_id: %s)", err, traceID)
	}

	response := StandardResponse{
		Success: false,
		Message: apiError.Message,
		Error:   &apiError,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Trace-ID", traceID)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendSuccess sends a successful response
func (eh *ErrorHandler) SendSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendCreated sends a created response
func (eh *ErrorHandler) SendCreated(w http.ResponseWriter, message string, data interface{}) {
	response := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// SendNoContent sends a no content response
func (eh *ErrorHandler) SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RecoveryMiddleware recovers from panics and handles them as errors
func (eh *ErrorHandler) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				eh.logger.Errorf("Panic recovered: %v", err)
				eh.HandleError(w, r, fmt.Errorf("internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ValidationMiddleware validates request content type and basic structure
func (eh *ErrorHandler) ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add trace ID if not present
		if r.Header.Get("X-Trace-ID") == "" {
			r.Header.Set("X-Trace-ID", uuid.New().String())
		}

		// Validate content type for POST/PUT requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" && contentType != "" {
				eh.HandleError(w, r, &ValidationError{
					Message: "Content-Type must be application/json",
					Fields:  map[string]string{"content-type": "Invalid content type"},
				})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
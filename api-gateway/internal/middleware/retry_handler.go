package middleware

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries int           // Maximum number of retries
	BaseDelay  time.Duration // Base delay between retries
	MaxDelay   time.Duration // Maximum delay between retries
	Multiplier float64       // Exponential backoff multiplier
	Jitter     bool          // Whether to add jitter to delays
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
		Jitter:     true,
	}
}

// RetryHandler handles retry logic with exponential backoff
type RetryHandler struct {
	config RetryConfig
	logger logger.Logger
}

// NewRetryHandler creates a new retry handler
func NewRetryHandler(config RetryConfig, logger logger.Logger) *RetryHandler {
	return &RetryHandler{
		config: config,
		logger: logger,
	}
}

// WithRetry executes an operation with retry logic
func (rh *RetryHandler) WithRetry(ctx context.Context, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= rh.config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lastErr = operation()
		if lastErr == nil {
			if attempt > 0 {
				rh.logger.Infof("Operation succeeded after %d attempts", attempt+1)
			}
			return nil
		}

		// Check if error is retryable
		if !rh.isRetryable(lastErr) {
			rh.logger.Warnf("Non-retryable error encountered: %v", lastErr)
			return lastErr
		}

		if attempt == rh.config.MaxRetries {
			rh.logger.Errorf("Operation failed after %d attempts: %v", attempt+1, lastErr)
			break
		}

		delay := rh.calculateDelay(attempt)
		rh.logger.Warnf("Operation failed (attempt %d/%d), retrying in %v: %v",
			attempt+1, rh.config.MaxRetries+1, delay, lastErr)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return lastErr
}

// calculateDelay calculates the delay for the given attempt
func (rh *RetryHandler) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(rh.config.BaseDelay) * math.Pow(rh.config.Multiplier, float64(attempt)))

	// Apply maximum delay limit
	if delay > rh.config.MaxDelay {
		delay = rh.config.MaxDelay
	}

	// Add jitter if enabled
	if rh.config.Jitter {
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1) // 10% jitter
		delay = delay + jitter
	}

	return delay
}

// isRetryable determines if an error is retryable
func (rh *RetryHandler) isRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for gRPC status codes
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
			return true
		case codes.Internal:
			// Some internal errors might be retryable
			return true
		default:
			return false
		}
	}

	// For non-gRPC errors, you can add custom logic here
	// For now, we'll be conservative and not retry
	return false
}

// RetryInterceptor creates a gRPC client interceptor with retry logic
func (rh *RetryHandler) RetryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return rh.WithRetry(ctx, func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}

// CombinedInterceptor combines circuit breaker and retry logic
func CombinedInterceptor(cb *CircuitBreaker, rh *RetryHandler) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return cb.Execute(func() error {
			return rh.WithRetry(ctx, func() error {
				// Apply timeout to the context
				timeoutCtx, cancel := context.WithTimeout(ctx, cb.timeout)
				defer cancel()

				return invoker(timeoutCtx, method, req, reply, cc, opts...)
			})
		})
	}
}

// HTTPRetryWrapper wraps HTTP operations with retry logic
func (rh *RetryHandler) HTTPRetryWrapper(ctx context.Context, operation func() error) error {
	return rh.WithRetry(ctx, operation)
}
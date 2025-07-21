package middleware

import (
	"context"
	"math"
	"time"

	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int           // Maximum number of retry attempts
	InitialDelay    time.Duration // Initial delay before first retry
	MaxDelay        time.Duration // Maximum delay between retries
	BackoffFactor   float64       // Exponential backoff factor
	RetryableErrors []codes.Code  // List of gRPC error codes that should be retried
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []codes.Code{
			codes.Unavailable,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
			codes.Aborted,
			codes.Internal,
		},
	}
}

// RetryInterceptor creates a gRPC interceptor with retry logic
func RetryInterceptor(config RetryConfig, logger logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var lastErr error
		
		for attempt := 0; attempt < config.MaxAttempts; attempt++ {
			// First attempt - no delay
			if attempt > 0 {
				delay := calculateDelay(attempt, config)
				logger.Infof("Retrying %s, attempt %d/%d after %v", method, attempt+1, config.MaxAttempts, delay)
				
				select {
				case <-time.After(delay):
					// Continue with retry
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			
			// Make the actual call
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				if attempt > 0 {
					logger.Infof("Retry successful for %s on attempt %d", method, attempt+1)
				}
				return nil
			}
			
			lastErr = err
			
			// Check if this error should be retried
			if !shouldRetry(err, config.RetryableErrors) {
				logger.Infof("Non-retryable error for %s: %v", method, err)
				return err
			}
			
			// Check if we should continue retrying
			if attempt+1 >= config.MaxAttempts {
				logger.Infof("Max retry attempts reached for %s: %v", method, err)
				break
			}
		}
		
		return lastErr
	}
}

// calculateDelay calculates the delay for a retry attempt using exponential backoff
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt-1)))
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return delay
}

// shouldRetry checks if an error should be retried based on the configuration
func shouldRetry(err error, retryableErrors []codes.Code) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	
	for _, code := range retryableErrors {
		if st.Code() == code {
			return true
		}
	}
	
	return false
}

// RetryManager manages retry configurations for different services
type RetryManager struct {
	configs map[string]RetryConfig
	logger  logger.Logger
}

// NewRetryManager creates a new retry manager
func NewRetryManager(logger logger.Logger) *RetryManager {
	return &RetryManager{
		configs: make(map[string]RetryConfig),
		logger:  logger,
	}
}

// SetConfig sets retry configuration for a service
func (rm *RetryManager) SetConfig(serviceName string, config RetryConfig) {
	rm.configs[serviceName] = config
}

// GetConfig gets retry configuration for a service
func (rm *RetryManager) GetConfig(serviceName string) RetryConfig {
	if config, exists := rm.configs[serviceName]; exists {
		return config
	}
	return DefaultRetryConfig()
}

// GetInterceptor returns a retry interceptor for a service
func (rm *RetryManager) GetInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	config := rm.GetConfig(serviceName)
	return RetryInterceptor(config, rm.logger)
}

// RetryWithBackoff executes a function with retry logic
func RetryWithBackoff(ctx context.Context, config RetryConfig, logger logger.Logger, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			delay := calculateDelay(attempt, config)
			logger.Infof("Retrying operation, attempt %d/%d after %v", attempt+1, config.MaxAttempts, delay)
			
			select {
			case <-time.After(delay):
				// Continue with retry
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		
		err := fn()
		if err == nil {
			if attempt > 0 {
				logger.Infof("Retry successful on attempt %d", attempt+1)
			}
			return nil
		}
		
		lastErr = err
		
		if attempt+1 >= config.MaxAttempts {
			logger.Infof("Max retry attempts reached: %v", err)
			break
		}
	}
	
	return lastErr
}
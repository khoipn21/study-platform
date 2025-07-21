package middleware

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	maxFailures   int
	timeout       time.Duration
	resetTimeout  time.Duration
	
	failures      int
	lastFailTime  time.Time
	state         CircuitBreakerState
	mutex         sync.RWMutex
	
	logger        logger.Logger
	serviceName   string
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures  int           // Maximum failures before opening
	Timeout      time.Duration // Timeout for requests
	ResetTimeout time.Duration // Time to wait before trying again
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(serviceName string, config CircuitBreakerConfig, logger logger.Logger) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  config.MaxFailures,
		timeout:      config.Timeout,
		resetTimeout: config.ResetTimeout,
		state:        StateClosed,
		logger:       logger,
		serviceName:  serviceName,
	}
}

// Execute runs the function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.canExecute() {
		return errors.New("circuit breaker is open")
	}
	
	err := fn()
	cb.recordResult(err)
	return err
}

// canExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = StateHalfOpen
			cb.mutex.Unlock()
			cb.mutex.RLock()
			cb.logger.Info("Circuit breaker for " + cb.serviceName + " moved to half-open state")
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult records the result of an execution
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.logger.Info("Circuit breaker for " + cb.serviceName + " opened from half-open state")
		} else if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			cb.logger.Info("Circuit breaker for " + cb.serviceName + " opened due to max failures")
		}
	} else {
		cb.failures = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			cb.logger.Info("Circuit breaker for " + cb.serviceName + " closed from half-open state")
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.failures
}

// CircuitBreakerInterceptor creates a gRPC interceptor with circuit breaker
func CircuitBreakerInterceptor(cb *CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return cb.Execute(func() error {
			// Apply timeout to the context
			timeoutCtx, cancel := context.WithTimeout(ctx, cb.timeout)
			defer cancel()
			
			return invoker(timeoutCtx, method, req, reply, cc, opts...)
		})
	}
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	config   CircuitBreakerConfig
	logger   logger.Logger
	mutex    sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(config CircuitBreakerConfig, logger logger.Logger) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
		logger:   logger,
	}
}

// GetCircuitBreaker gets or creates a circuit breaker for a service
func (cbm *CircuitBreakerManager) GetCircuitBreaker(serviceName string) *CircuitBreaker {
	cbm.mutex.RLock()
	if cb, exists := cbm.breakers[serviceName]; exists {
		cbm.mutex.RUnlock()
		return cb
	}
	cbm.mutex.RUnlock()
	
	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()
	
	// Double-check pattern
	if cb, exists := cbm.breakers[serviceName]; exists {
		return cb
	}
	
	cb := NewCircuitBreaker(serviceName, cbm.config, cbm.logger)
	cbm.breakers[serviceName] = cb
	return cb
}

// GetStatus returns the status of all circuit breakers
func (cbm *CircuitBreakerManager) GetStatus() map[string]interface{} {
	cbm.mutex.RLock()
	defer cbm.mutex.RUnlock()
	
	status := make(map[string]interface{})
	for name, cb := range cbm.breakers {
		state := cb.GetState()
		var stateStr string
		switch state {
		case StateClosed:
			stateStr = "closed"
		case StateOpen:
			stateStr = "open"
		case StateHalfOpen:
			stateStr = "half-open"
		}
		
		status[name] = map[string]interface{}{
			"state":    stateStr,
			"failures": cb.GetFailureCount(),
		}
	}
	
	return status
}
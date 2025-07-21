package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/study-platform/pkg/logger"
)

// RateLimiter represents a token bucket rate limiter
type RateLimiter struct {
	rate     int           // tokens per second
	capacity int           // bucket capacity
	tokens   int           // current tokens
	lastTime time.Time     // last refill time
	mutex    sync.Mutex    // thread safety
}

// RateLimitMiddleware manages rate limiting for the API Gateway
type RateLimitMiddleware struct {
	limiters map[string]*RateLimiter // IP-based rate limiters
	mutex    sync.RWMutex            // thread safety for map
	logger   logger.Logger
	// Global rate limiting settings
	globalRate     int // requests per second
	globalCapacity int // burst capacity
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(logger logger.Logger, rate, capacity int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiters:       make(map[string]*RateLimiter),
		logger:         logger,
		globalRate:     rate,
		globalCapacity: capacity,
	}
}

// newRateLimiter creates a new token bucket rate limiter
func newRateLimiter(rate, capacity int) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: capacity,
		tokens:   capacity,
		lastTime: time.Now(),
	}
}

// allow checks if a request is allowed based on rate limiting
func (rl *RateLimiter) allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(rl.lastTime)
	tokensToAdd := int(elapsed.Seconds() * float64(rl.rate))
	
	// Add tokens but don't exceed capacity
	rl.tokens += tokensToAdd
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	
	rl.lastTime = now

	// Check if we have tokens available
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	
	return false
}

// getLimiterForIP gets or creates a rate limiter for an IP address
func (rlm *RateLimitMiddleware) getLimiterForIP(ip string) *RateLimiter {
	rlm.mutex.RLock()
	limiter, exists := rlm.limiters[ip]
	rlm.mutex.RUnlock()
	
	if exists {
		return limiter
	}
	
	// Create new limiter
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()
	
	// Double-check pattern
	if limiter, exists := rlm.limiters[ip]; exists {
		return limiter
	}
	
	limiter = newRateLimiter(rlm.globalRate, rlm.globalCapacity)
	rlm.limiters[ip] = limiter
	return limiter
}

// getClientIP extracts the client IP from the request
func (rlm *RateLimitMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// RateLimit is the middleware function that applies rate limiting
func (rlm *RateLimitMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := rlm.getClientIP(r)
		limiter := rlm.getLimiterForIP(clientIP)
		
		if !limiter.allow() {
			rlm.logger.Info("Rate limit exceeded for IP: " + clientIP)
			
			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", "60")
			
			// Send rate limit response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			
			response := `{"success":false,"message":"Rate limit exceeded. Too many requests.","error":"rate_limit_exceeded"}`
			w.Write([]byte(response))
			return
		}
		
		// Set rate limit headers for successful requests
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "99") // Simplified
		
		next.ServeHTTP(w, r)
	})
}

// cleanup removes old rate limiters to prevent memory leaks
func (rlm *RateLimitMiddleware) cleanup() {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()
	
	// Remove limiters that haven't been used in the last hour
	cutoff := time.Now().Add(-time.Hour)
	for ip, limiter := range rlm.limiters {
		if limiter.lastTime.Before(cutoff) {
			delete(rlm.limiters, ip)
		}
	}
}

// StartCleanupRoutine starts a goroutine to periodically clean up old rate limiters
func (rlm *RateLimitMiddleware) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			rlm.cleanup()
		}
	}()
}
package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/study-platform/pkg/logger"
)

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status    string                    `json:"status"`
	Timestamp time.Time                 `json:"timestamp"`
	Services  map[string]ServiceHealth  `json:"services"`
	Version   string                    `json:"version"`
	Uptime    string                    `json:"uptime"`
}

// ServiceHealth represents individual service health
type ServiceHealth struct {
	Status    string        `json:"status"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	LastCheck time.Time     `json:"last_check"`
}

// ExternalService represents a service to check
type ExternalService struct {
	Name     string
	URL      string
	Timeout  time.Duration
	Critical bool // If true, service failure affects overall status
}

// HealthChecker manages health checks
type HealthChecker struct {
	services  map[string]ExternalService
	cache     map[string]ServiceHealth
	mutex     sync.RWMutex
	logger    logger.Logger
	startTime time.Time
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(logger logger.Logger) *HealthChecker {
	return &HealthChecker{
		services:  make(map[string]ExternalService),
		cache:     make(map[string]ServiceHealth),
		logger:    logger,
		startTime: time.Now(),
	}
}

// RegisterService registers a service for health checking
func (hc *HealthChecker) RegisterService(name, url string, timeout time.Duration, critical bool) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.services[name] = ExternalService{
		Name:     name,
		URL:      url,
		Timeout:  timeout,
		Critical: critical,
	}
}

// CheckHealth performs health checks on all registered services
func (hc *HealthChecker) CheckHealth(ctx context.Context) *HealthStatus {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]ServiceHealth),
		Version:   getVersion(),
		Uptime:    time.Since(hc.startTime).String(),
	}

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		name   string
		health ServiceHealth
	}, len(hc.services))

	// Check all services concurrently
	for name, service := range hc.services {
		wg.Add(1)
		go func(name string, service ExternalService) {
			defer wg.Done()
			health := hc.checkSingleService(ctx, service)
			resultChan <- struct {
				name   string
				health ServiceHealth
			}{name, health}
		}(name, service)
	}

	// Wait for all checks to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		status.Services[result.name] = result.health
		hc.cache[result.name] = result.health

		// If any critical service is unhealthy, mark overall status as unhealthy
		service := hc.services[result.name]
		if service.Critical && result.health.Status != "healthy" {
			status.Status = "unhealthy"
		}
	}

	return status
}

// checkSingleService checks the health of a single service
func (hc *HealthChecker) checkSingleService(ctx context.Context, service ExternalService) ServiceHealth {
	start := time.Now()
	health := ServiceHealth{
		LastCheck: start,
		Latency:   0,
	}

	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Make HTTP request
	req, err := http.NewRequestWithContext(checkCtx, "GET", service.URL, nil)
	if err != nil {
		health.Status = "unhealthy"
		health.Error = fmt.Sprintf("Failed to create request: %v", err)
		return health
	}

	client := &http.Client{
		Timeout: service.Timeout,
	}

	resp, err := client.Do(req)
	health.Latency = time.Since(start)

	if err != nil {
		health.Status = "unhealthy"
		health.Error = fmt.Sprintf("Request failed: %v", err)
		return health
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		health.Status = "healthy"
	} else {
		health.Status = "unhealthy"
		health.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return health
}

// GetCachedHealth returns cached health status
func (hc *HealthChecker) GetCachedHealth() *HealthStatus {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]ServiceHealth),
		Version:   getVersion(),
		Uptime:    time.Since(hc.startTime).String(),
	}

	for name, health := range hc.cache {
		status.Services[name] = health

		// Check if service is critical and unhealthy
		if service, exists := hc.services[name]; exists && service.Critical && health.Status != "healthy" {
			status.Status = "unhealthy"
		}
	}

	return status
}

// HealthHandler handles HTTP health check requests
func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	var health *HealthStatus

	// Use cached health for non-detailed checks
	if r.URL.Query().Get("detailed") == "true" {
		health = hc.CheckHealth(r.Context())
	} else {
		health = hc.GetCachedHealth()
	}

	w.Header().Set("Content-Type", "application/json")

	// Set HTTP status based on health
	if health.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(health)
}

// ReadinessHandler handles readiness probe requests
func (hc *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	health := hc.GetCachedHealth()

	// Check only critical services for readiness
	ready := true
	for name, serviceHealth := range health.Services {
		if service, exists := hc.services[name]; exists && service.Critical {
			if serviceHealth.Status != "healthy" {
				ready = false
				break
			}
		}
	}

	response := map[string]interface{}{
		"ready": ready,
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// LivenessHandler handles liveness probe requests
func (hc *HealthChecker) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"alive": true,
		"timestamp": time.Now(),
		"uptime": time.Since(hc.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getVersion returns the service version from environment variable
func getVersion() string {
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		return "unknown"
	}
	return version
}
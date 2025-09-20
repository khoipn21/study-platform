package handler

import (
	"context"
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := h.getSystemStatus()

	httpStatus := http.StatusOK
	if status["status"] != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, status)
}

// ReadinessCheck handles GET /health/ready
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ready := h.isReady()

	httpStatus := http.StatusOK
	if !ready {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": map[string]bool{"ready": ready},
		"service": "instructor-dashboard",
		"timestamp": time.Now().UTC(),
	})
}

// LivenessCheck handles GET /health/live
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Simple liveness check - service is alive if it can respond
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"service": "instructor-dashboard",
		"timestamp": time.Now().UTC(),
	})
}

// MetricsEndpoint handles GET /metrics
func (h *HealthHandler) MetricsEndpoint(c *gin.Context) {
	metrics := h.collectMetrics()
	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"service": "instructor-dashboard",
		"timestamp": time.Now().UTC(),
	})
}

// Helper methods

func (h *HealthHandler) getSystemStatus() gin.H {
	status := "healthy"
	checks := make(map[string]interface{})

	// Database connectivity check
	dbStatus := h.checkDatabase()
	checks["database"] = dbStatus
	if !dbStatus["healthy"].(bool) {
		status = "unhealthy"
	}

	// Memory check (basic)
	memStatus := h.checkMemory()
	checks["memory"] = memStatus
	if !memStatus["healthy"].(bool) {
		status = "degraded"
	}

	return gin.H{
		"status": status,
		"service": "instructor-dashboard",
		"version": "1.0.0",
		"checks": checks,
		"timestamp": time.Now().UTC(),
	}
}

func (h *HealthHandler) isReady() bool {
	// Check if database is accessible
	if !h.checkDatabase()["healthy"].(bool) {
		return false
	}

	// Add other readiness checks here
	// e.g., external service dependencies

	return true
}

func (h *HealthHandler) checkDatabase() map[string]interface{} {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)
	duration := time.Since(start)

	if err != nil {
		return map[string]interface{}{
			"healthy": false,
			"error": err.Error(),
			"response_time_ms": duration.Milliseconds(),
		}
	}

	return map[string]interface{}{
		"healthy": true,
		"response_time_ms": duration.Milliseconds(),
	}
}

func (h *HealthHandler) checkMemory() map[string]interface{} {
	// Basic memory check - could be enhanced with more sophisticated monitoring
	return map[string]interface{}{
		"healthy": true,
		"status": "ok",
	}
}

func (h *HealthHandler) collectMetrics() map[string]interface{} {
	return map[string]interface{}{
		"uptime_seconds": time.Since(startTime).Seconds(),
		"database_connections": h.getDatabaseConnections(),
		"goroutines": runtime.NumGoroutine(),
		"memory": h.getMemoryStats(),
	}
}

func (h *HealthHandler) getDatabaseConnections() map[string]interface{} {
	stats := h.db.Stats()
	return map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use": stats.InUse,
		"idle": stats.Idle,
		"wait_count": stats.WaitCount,
		"wait_duration_ms": stats.WaitDuration.Milliseconds(),
		"max_idle_closed": stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}

func (h *HealthHandler) getMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"alloc_bytes": m.Alloc,
		"total_alloc_bytes": m.TotalAlloc,
		"sys_bytes": m.Sys,
		"num_gc": m.NumGC,
	}
}

var startTime = time.Now()
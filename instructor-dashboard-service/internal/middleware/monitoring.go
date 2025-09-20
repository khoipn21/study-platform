package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"instructor-dashboard-service/internal/monitoring"
)

// MonitoringMiddleware provides comprehensive request monitoring
func MonitoringMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Writer.Status()

		// Record metrics
		endpoint := method + " " + path
		hasError := status >= 400

		monitoring.RecordRequest(endpoint, duration, hasError)

		// Add monitoring headers
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Request-ID", c.GetString("request_id"))
	})
}

// RequestIDMiddleware adds unique request IDs for tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	})
}

// MetricsEndpoint provides metrics endpoint
func MetricsEndpoint(collector *monitoring.MetricsCollector) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		format := c.Query("format")

		switch format {
		case "prometheus":
			c.Header("Content-Type", "text/plain")
			c.String(200, collector.ExportPrometheusMetrics())
		default:
			metricsJSON, err := collector.ExportMetrics()
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to export metrics"})
				return
			}
			c.Header("Content-Type", "application/json")
			c.Data(200, "application/json", metricsJSON)
		}
	})
}

// AlertsEndpoint provides alerts endpoint
func AlertsEndpoint(collector *monitoring.MetricsCollector) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		alerts := collector.CheckAlerts()
		c.JSON(200, gin.H{
			"alerts": alerts,
			"total":  len(alerts),
		})
	})
}

func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
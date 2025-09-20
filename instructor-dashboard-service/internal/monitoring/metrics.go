package monitoring

import (
	"context"
	"database/sql"
	"encoding/json"
	"runtime"
	"sync"
	"time"
)

// MetricsCollector collects and aggregates metrics for the instructor dashboard
type MetricsCollector struct {
	db           *sql.DB
	metrics      map[string]interface{}
	mutex        sync.RWMutex
	lastCollected time.Time
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	ServiceName    string            `json:"service_name"`
	Version        string            `json:"version"`
	Uptime         int64             `json:"uptime_seconds"`
	Timestamp      time.Time         `json:"timestamp"`
	Memory         MemoryMetrics     `json:"memory"`
	Database       DatabaseMetrics   `json:"database"`
	Performance    PerformanceMetrics `json:"performance"`
	Business       BusinessMetrics   `json:"business"`
	Errors         ErrorMetrics      `json:"errors"`
	CustomMetrics  map[string]interface{} `json:"custom_metrics"`
}

type MemoryMetrics struct {
	AllocBytes      uint64 `json:"alloc_bytes"`
	TotalAllocBytes uint64 `json:"total_alloc_bytes"`
	SysBytes        uint64 `json:"sys_bytes"`
	HeapObjects     uint64 `json:"heap_objects"`
	GCCycles        uint32 `json:"gc_cycles"`
	Goroutines      int    `json:"goroutines"`
}

type DatabaseMetrics struct {
	OpenConnections     int           `json:"open_connections"`
	InUse              int           `json:"in_use"`
	Idle               int           `json:"idle"`
	WaitCount          int64         `json:"wait_count"`
	WaitDuration       time.Duration `json:"wait_duration_ms"`
	MaxIdleClosed      int64         `json:"max_idle_closed"`
	MaxLifetimeClosed  int64         `json:"max_lifetime_closed"`
	QueryLatencyP95    float64       `json:"query_latency_p95_ms"`
	ActiveQueries      int           `json:"active_queries"`
}

type PerformanceMetrics struct {
	RequestsTotal         int64   `json:"requests_total"`
	RequestsPerSecond     float64 `json:"requests_per_second"`
	AverageResponseTime   float64 `json:"avg_response_time_ms"`
	P95ResponseTime       float64 `json:"p95_response_time_ms"`
	P99ResponseTime       float64 `json:"p99_response_time_ms"`
	ActiveConnections     int     `json:"active_connections"`
	TotalConnections      int64   `json:"total_connections"`
}

type BusinessMetrics struct {
	ActiveInstructors     int     `json:"active_instructors"`
	TotalDashboardViews   int64   `json:"total_dashboard_views"`
	CoursesCreatedToday   int     `json:"courses_created_today"`
	RevenueToday          float64 `json:"revenue_today"`
	AnalyticsRequests     int64   `json:"analytics_requests"`
	BroadcastMessages     int64   `json:"broadcast_messages"`
	FileUploads           int64   `json:"file_uploads"`
}

type ErrorMetrics struct {
	ErrorsTotal         int64            `json:"errors_total"`
	ErrorRate           float64          `json:"error_rate"`
	ErrorsByType        map[string]int64 `json:"errors_by_type"`
	ErrorsByEndpoint    map[string]int64 `json:"errors_by_endpoint"`
	Last5MinutesErrors  int64            `json:"last_5min_errors"`
	CriticalErrors      int64            `json:"critical_errors"`
}

var (
	startTime = time.Now()
	requestMetrics = make(map[string]*RequestMetrics)
	metricsLock sync.RWMutex
)

type RequestMetrics struct {
	Count         int64
	TotalDuration time.Duration
	LastRequest   time.Time
	Errors        int64
}

func NewMetricsCollector(db *sql.DB) *MetricsCollector {
	return &MetricsCollector{
		db:      db,
		metrics: make(map[string]interface{}),
	}
}

// CollectSystemMetrics collects comprehensive system metrics
func (mc *MetricsCollector) CollectSystemMetrics() *SystemMetrics {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	now := time.Now()

	metrics := &SystemMetrics{
		ServiceName:   "instructor-dashboard-service",
		Version:       "1.0.0",
		Uptime:        int64(now.Sub(startTime).Seconds()),
		Timestamp:     now,
		Memory:        mc.collectMemoryMetrics(),
		Database:      mc.collectDatabaseMetrics(),
		Performance:   mc.collectPerformanceMetrics(),
		Business:      mc.collectBusinessMetrics(),
		Errors:        mc.collectErrorMetrics(),
		CustomMetrics: make(map[string]interface{}),
	}

	mc.lastCollected = now
	return metrics
}

// CollectMemoryMetrics collects memory and runtime metrics
func (mc *MetricsCollector) collectMemoryMetrics() MemoryMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryMetrics{
		AllocBytes:      m.Alloc,
		TotalAllocBytes: m.TotalAlloc,
		SysBytes:        m.Sys,
		HeapObjects:     m.HeapObjects,
		GCCycles:        m.NumGC,
		Goroutines:      runtime.NumGoroutine(),
	}
}

// collectDatabaseMetrics collects database performance metrics
func (mc *MetricsCollector) collectDatabaseMetrics() DatabaseMetrics {
	stats := mc.db.Stats()

	// Sample query to measure latency
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := mc.db.QueryContext(ctx, "SELECT 1")
	queryLatency := time.Since(start).Seconds() * 1000 // Convert to milliseconds

	if err != nil {
		queryLatency = -1 // Indicate error
	}

	return DatabaseMetrics{
		OpenConnections:    stats.OpenConnections,
		InUse:             stats.InUse,
		Idle:              stats.Idle,
		WaitCount:         stats.WaitCount,
		WaitDuration:      stats.WaitDuration,
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,
		QueryLatencyP95:   queryLatency,
		ActiveQueries:     stats.InUse, // Approximation
	}
}

// collectPerformanceMetrics collects API performance metrics
func (mc *MetricsCollector) collectPerformanceMetrics() PerformanceMetrics {
	metricsLock.RLock()
	defer metricsLock.RUnlock()

	var totalRequests int64
	var totalDuration time.Duration
	var totalErrors int64

	for _, rm := range requestMetrics {
		totalRequests += rm.Count
		totalDuration += rm.TotalDuration
		totalErrors += rm.Errors
	}

	avgResponseTime := float64(0)
	if totalRequests > 0 {
		avgResponseTime = float64(totalDuration.Nanoseconds()) / float64(totalRequests) / 1000000 // Convert to ms
	}

	uptime := time.Since(startTime).Seconds()
	rps := float64(totalRequests) / uptime

	return PerformanceMetrics{
		RequestsTotal:       totalRequests,
		RequestsPerSecond:   rps,
		AverageResponseTime: avgResponseTime,
		P95ResponseTime:     mc.calculateP95ResponseTime(),
		P99ResponseTime:     mc.calculateP99ResponseTime(),
		ActiveConnections:   runtime.NumGoroutine(), // Approximation
		TotalConnections:    totalRequests,
	}
}

// collectBusinessMetrics collects business-specific metrics
func (mc *MetricsCollector) collectBusinessMetrics() BusinessMetrics {
	metrics := BusinessMetrics{}

	// Collect active instructors (last 24 hours)
	row := mc.db.QueryRow(`
		SELECT COUNT(DISTINCT instructor_id)
		FROM instructor_sessions
		WHERE last_activity > NOW() - INTERVAL '24 hours'
	`)
	row.Scan(&metrics.ActiveInstructors)

	// Dashboard views today
	row = mc.db.QueryRow(`
		SELECT COUNT(*)
		FROM dashboard_views
		WHERE DATE(viewed_at) = CURRENT_DATE
	`)
	row.Scan(&metrics.TotalDashboardViews)

	// Courses created today
	row = mc.db.QueryRow(`
		SELECT COUNT(*)
		FROM courses
		WHERE DATE(created_at) = CURRENT_DATE
	`)
	row.Scan(&metrics.CoursesCreatedToday)

	// Revenue today
	row = mc.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM instructor_earnings
		WHERE DATE(earned_at) = CURRENT_DATE
	`)
	row.Scan(&metrics.RevenueToday)

	// Analytics requests (from cache or recent table)
	metricsLock.RLock()
	if rm, exists := requestMetrics["/api/v1/analytics"]; exists {
		metrics.AnalyticsRequests = rm.Count
	}
	metricsLock.RUnlock()

	return metrics
}

// collectErrorMetrics collects error metrics
func (mc *MetricsCollector) collectErrorMetrics() ErrorMetrics {
	metricsLock.RLock()
	defer metricsLock.RUnlock()

	var totalErrors int64
	var totalRequests int64
	errorsByType := make(map[string]int64)
	errorsByEndpoint := make(map[string]int64)

	for endpoint, rm := range requestMetrics {
		totalRequests += rm.Count
		totalErrors += rm.Errors
		if rm.Errors > 0 {
			errorsByEndpoint[endpoint] = rm.Errors
		}
	}

	errorRate := float64(0)
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}

	return ErrorMetrics{
		ErrorsTotal:        totalErrors,
		ErrorRate:          errorRate,
		ErrorsByType:       errorsByType,
		ErrorsByEndpoint:   errorsByEndpoint,
		Last5MinutesErrors: mc.getRecentErrors(5 * time.Minute),
		CriticalErrors:     mc.getCriticalErrors(),
	}
}

// RecordRequest records metrics for an HTTP request
func RecordRequest(endpoint string, duration time.Duration, hasError bool) {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	if requestMetrics[endpoint] == nil {
		requestMetrics[endpoint] = &RequestMetrics{}
	}

	rm := requestMetrics[endpoint]
	rm.Count++
	rm.TotalDuration += duration
	rm.LastRequest = time.Now()

	if hasError {
		rm.Errors++
	}
}

// ExportMetrics exports metrics in JSON format
func (mc *MetricsCollector) ExportMetrics() ([]byte, error) {
	metrics := mc.CollectSystemMetrics()
	return json.MarshalIndent(metrics, "", "  ")
}

// ExportPrometheusMetrics exports metrics in Prometheus format
func (mc *MetricsCollector) ExportPrometheusMetrics() string {
	metrics := mc.CollectSystemMetrics()

	prometheus := `# HELP instructor_dashboard_uptime_seconds Total uptime of the service
# TYPE instructor_dashboard_uptime_seconds counter
instructor_dashboard_uptime_seconds ` + toString(metrics.Uptime) + `

# HELP instructor_dashboard_memory_alloc_bytes Current allocated memory in bytes
# TYPE instructor_dashboard_memory_alloc_bytes gauge
instructor_dashboard_memory_alloc_bytes ` + toString(metrics.Memory.AllocBytes) + `

# HELP instructor_dashboard_goroutines Current number of goroutines
# TYPE instructor_dashboard_goroutines gauge
instructor_dashboard_goroutines ` + toString(metrics.Memory.Goroutines) + `

# HELP instructor_dashboard_db_connections_open Current number of open database connections
# TYPE instructor_dashboard_db_connections_open gauge
instructor_dashboard_db_connections_open ` + toString(metrics.Database.OpenConnections) + `

# HELP instructor_dashboard_requests_total Total number of HTTP requests
# TYPE instructor_dashboard_requests_total counter
instructor_dashboard_requests_total ` + toString(metrics.Performance.RequestsTotal) + `

# HELP instructor_dashboard_requests_per_second Current requests per second
# TYPE instructor_dashboard_requests_per_second gauge
instructor_dashboard_requests_per_second ` + toString(metrics.Performance.RequestsPerSecond) + `

# HELP instructor_dashboard_response_time_avg Average response time in milliseconds
# TYPE instructor_dashboard_response_time_avg gauge
instructor_dashboard_response_time_avg ` + toString(metrics.Performance.AverageResponseTime) + `

# HELP instructor_dashboard_active_instructors Number of active instructors in last 24h
# TYPE instructor_dashboard_active_instructors gauge
instructor_dashboard_active_instructors ` + toString(metrics.Business.ActiveInstructors) + `

# HELP instructor_dashboard_revenue_today Revenue generated today
# TYPE instructor_dashboard_revenue_today gauge
instructor_dashboard_revenue_today ` + toString(metrics.Business.RevenueToday) + `

# HELP instructor_dashboard_errors_total Total number of errors
# TYPE instructor_dashboard_errors_total counter
instructor_dashboard_errors_total ` + toString(metrics.Errors.ErrorsTotal) + `

# HELP instructor_dashboard_error_rate Error rate percentage
# TYPE instructor_dashboard_error_rate gauge
instructor_dashboard_error_rate ` + toString(metrics.Errors.ErrorRate) + `
`

	return prometheus
}

// Helper methods

func (mc *MetricsCollector) calculateP95ResponseTime() float64 {
	// Simplified implementation - in production, use proper percentile calculation
	metricsLock.RLock()
	defer metricsLock.RUnlock()

	if len(requestMetrics) == 0 {
		return 0
	}

	var totalDuration time.Duration
	var totalRequests int64

	for _, rm := range requestMetrics {
		totalDuration += rm.TotalDuration
		totalRequests += rm.Count
	}

	if totalRequests == 0 {
		return 0
	}

	avgMs := float64(totalDuration.Nanoseconds()) / float64(totalRequests) / 1000000
	return avgMs * 1.2 // Rough P95 approximation
}

func (mc *MetricsCollector) calculateP99ResponseTime() float64 {
	p95 := mc.calculateP95ResponseTime()
	return p95 * 1.3 // Rough P99 approximation
}

func (mc *MetricsCollector) getRecentErrors(duration time.Duration) int64 {
	// In production, this would query an errors table with timestamps
	return 0
}

func (mc *MetricsCollector) getCriticalErrors() int64 {
	// In production, this would query critical errors from logs/database
	return 0
}

// AlertingRule defines rules for metric-based alerting
type AlertingRule struct {
	Name        string      `json:"name"`
	Metric      string      `json:"metric"`
	Operator    string      `json:"operator"` // "gt", "lt", "eq", "gte", "lte"
	Threshold   interface{} `json:"threshold"`
	Duration    string      `json:"duration"`
	Severity    string      `json:"severity"` // "low", "medium", "high", "critical"
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
}

var defaultAlertingRules = []AlertingRule{
	{
		Name:        "high_error_rate",
		Metric:      "error_rate",
		Operator:    "gt",
		Threshold:   5.0,
		Duration:    "5m",
		Severity:    "high",
		Description: "Error rate exceeds 5% for 5 minutes",
		Enabled:     true,
	},
	{
		Name:        "high_response_time",
		Metric:      "avg_response_time",
		Operator:    "gt",
		Threshold:   3000.0,
		Duration:    "5m",
		Severity:    "medium",
		Description: "Average response time exceeds 3 seconds",
		Enabled:     true,
	},
	{
		Name:        "database_connection_exhaustion",
		Metric:      "db_connections_in_use",
		Operator:    "gt",
		Threshold:   80,
		Duration:    "2m",
		Severity:    "critical",
		Description: "Database connections usage above 80%",
		Enabled:     true,
	},
}

// CheckAlerts evaluates alerting rules against current metrics
func (mc *MetricsCollector) CheckAlerts() []AlertingRule {
	triggeredAlerts := []AlertingRule{}
	metrics := mc.CollectSystemMetrics()

	for _, rule := range defaultAlertingRules {
		if !rule.Enabled {
			continue
		}

		if mc.evaluateRule(rule, metrics) {
			triggeredAlerts = append(triggeredAlerts, rule)
		}
	}

	return triggeredAlerts
}

func (mc *MetricsCollector) evaluateRule(rule AlertingRule, metrics *SystemMetrics) bool {
	var value float64

	// Extract metric value based on rule.Metric
	switch rule.Metric {
	case "error_rate":
		value = metrics.Errors.ErrorRate
	case "avg_response_time":
		value = metrics.Performance.AverageResponseTime
	case "db_connections_in_use":
		value = float64(metrics.Database.InUse)
	default:
		return false
	}

	// Evaluate condition
	threshold, ok := rule.Threshold.(float64)
	if !ok {
		return false
	}

	switch rule.Operator {
	case "gt":
		return value > threshold
	case "gte":
		return value >= threshold
	case "lt":
		return value < threshold
	case "lte":
		return value <= threshold
	case "eq":
		return value == threshold
	default:
		return false
	}
}

// Utility function to convert interface{} to string for Prometheus
func toString(v interface{}) string {
	switch val := v.(type) {
	case int:
		return string(rune(val))
	case int64:
		return string(rune(val))
	case uint64:
		return string(rune(val))
	case float64:
		return string(rune(val))
	default:
		return "0"
	}
}
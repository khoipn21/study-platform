package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Service   string    `json:"service"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		Service:   "bucket-service",
		Database:  "connected",
	}

	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status.Status = "unhealthy"
		status.Database = "disconnected"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// More thorough checks for readiness
	status := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		Service:   "bucket-service",
		Database:  "connected",
	}

	// Check database connection and perform a simple query
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status.Status = "not ready"
		status.Database = "disconnected"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	// Try to perform a simple database operation
	var count int64
	if err := h.db.Table("files").Count(&count).Error; err != nil {
		status.Status = "not ready"
		status.Database = "query failed"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
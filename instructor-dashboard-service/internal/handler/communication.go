package handler

import (
	"net/http"
	"time"

	"instructor-dashboard-service/internal/model"
	"instructor-dashboard-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommunicationHandler struct {
	communicationService *service.CommunicationService
}

func NewCommunicationHandler(communicationService *service.CommunicationService) *CommunicationHandler {
	return &CommunicationHandler{
		communicationService: communicationService,
	}
}

// SendBroadcast handles POST /api/v1/instructor/communication/broadcast
func (h *CommunicationHandler) SendBroadcast(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var req model.SendBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.communicationService.SendBroadcast(instructorID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send broadcast message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Broadcast message sent successfully",
	})
}

// GetCommunicationHistory handles GET /api/v1/instructor/communication/history
func (h *CommunicationHandler) GetCommunicationHistory(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	// Parse pagination parameters
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)

	pagination := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse filter parameters
	filter := &model.Filter{
		Status: c.Query("status"),
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			filter.CourseID = &courseID
		}
	}

	communications, paginationResult, err := h.communicationService.GetCommunicationHistory(instructorID, filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get communication history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       communications,
		"pagination": paginationResult,
	})
}

// SetupAutomatedMessages handles POST /api/v1/instructor/communication/automated
func (h *CommunicationHandler) SetupAutomatedMessages(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var rules []model.AutomatedMessageRule
	if err := c.ShouldBindJSON(&rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.communicationService.SetupAutomatedMessages(instructorID, rules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup automated messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Automated messages setup successfully",
	})
}

// GetAutomatedMessages handles GET /api/v1/instructor/communication/automated
func (h *CommunicationHandler) GetAutomatedMessages(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	rules, err := h.communicationService.GetAutomatedMessageRules(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get automated message rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
	})
}

// SendPersonalMessage handles POST /api/v1/instructor/communication/personal/:student_id
func (h *CommunicationHandler) SendPersonalMessage(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	studentIDStr := c.Param("student_id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var req model.PersonalMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.communicationService.SendPersonalMessage(instructorID, studentID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send personal message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Personal message sent successfully",
	})
}

// MarkMessageRead handles PUT /api/v1/instructor/communication/:id/read
func (h *CommunicationHandler) MarkMessageRead(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	communicationIDStr := c.Param("id")
	communicationID, err := uuid.Parse(communicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	err = h.communicationService.MarkMessageRead(instructorID, communicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message marked as read",
	})
}

// UpdateCommunicationStatus handles PUT /api/v1/instructor/communication/:id/status
func (h *CommunicationHandler) UpdateCommunicationStatus(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	communicationIDStr := c.Param("id")
	communicationID, err := uuid.Parse(communicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.communicationService.UpdateCommunicationStatus(instructorID, communicationID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update communication status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Communication status updated successfully",
	})
}

// GetMessageTemplates handles GET /api/v1/instructor/communication/templates
func (h *CommunicationHandler) GetMessageTemplates(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	templates, err := h.communicationService.GetCommunicationTemplates(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    templates,
	})
}

// GetNotifications handles GET /api/v1/instructor/notifications
func (h *CommunicationHandler) GetNotifications(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	// Parse limit parameter
	limit := parseIntParam(c, "limit", 10)
	if limit > 50 {
		limit = 50 // Cap at 50 for performance
	}

	// Parse status filter
	status := c.Query("status") // "read", "unread", or empty for all

	// Get notifications
	notifications, err := h.communicationService.GetNotifications(instructorID, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifications,
	})
}
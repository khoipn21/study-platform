package handler

import (
	"net/http"
	"strconv"
	"time"

	"instructor-dashboard-service/internal/model"
	"instructor-dashboard-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardOverview handles GET /api/v1/instructor/dashboard/overview
func (h *DashboardHandler) GetDashboardOverview(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	overview, err := h.dashboardService.GetDashboardOverview(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard overview"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overview,
	})
}

// UpdateDashboardSettings handles PUT /api/v1/instructor/dashboard/settings
func (h *DashboardHandler) UpdateDashboardSettings(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var req model.UpdateDashboardSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.dashboardService.UpdateDashboardSettings(instructorID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dashboard settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard settings updated successfully",
	})
}

// GetInstructorCourses handles GET /api/v1/instructor/courses
func (h *DashboardHandler) GetInstructorCourses(c *gin.Context) {
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

	courses, paginationResult, err := h.dashboardService.GetInstructorCourses(instructorID, filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get instructor courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       courses,
		"pagination": paginationResult,
	})
}

// BulkCourseOperations handles POST /api/v1/instructor/courses/:id/bulk-operations
func (h *DashboardHandler) BulkCourseOperations(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var req model.BulkCourseOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.dashboardService.BulkCourseOperations(instructorID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform bulk operations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bulk operations completed successfully",
	})
}

// GetStudents handles GET /api/v1/instructor/students
func (h *DashboardHandler) GetStudents(c *gin.Context) {
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

	filter := &model.Filter{}

	students, paginationResult, err := h.dashboardService.GetStudents(instructorID, filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       students,
		"pagination": paginationResult,
	})
}

// GetStudentDetails handles GET /api/v1/instructor/students/:id
func (h *DashboardHandler) GetStudentDetails(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	studentIDStr := c.Param("id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.dashboardService.GetStudentDetails(instructorID, studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
	})
}

// GetAISuggestions handles GET /api/v1/instructor/suggestions
func (h *DashboardHandler) GetAISuggestions(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	limit := parseIntParam(c, "limit", 10)

	suggestions, err := h.dashboardService.GetAISuggestions(instructorID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI suggestions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    suggestions,
	})
}

// ImplementSuggestion handles POST /api/v1/instructor/suggestions/:id/implement
func (h *DashboardHandler) ImplementSuggestion(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	suggestionIDStr := c.Param("id")
	suggestionID, err := uuid.Parse(suggestionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid suggestion ID"})
		return
	}

	err = h.dashboardService.ImplementSuggestion(instructorID, suggestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to implement suggestion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Suggestion implemented successfully",
	})
}

// GetTeamMembers handles GET /api/v1/instructor/team
func (h *DashboardHandler) GetTeamMembers(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	teamMembers, err := h.dashboardService.GetTeamMembers(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    teamMembers,
	})
}

// InviteTeamMember handles POST /api/v1/instructor/team/invite
func (h *DashboardHandler) InviteTeamMember(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var req model.InviteTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.dashboardService.InviteTeamMember(instructorID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite team member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team member invitation sent successfully",
	})
}

// UpdateTeamMember handles PUT /api/v1/instructor/team/:id
func (h *DashboardHandler) UpdateTeamMember(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	teamMemberIDStr := c.Param("id")
	teamMemberID, err := uuid.Parse(teamMemberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team member ID"})
		return
	}

	var req model.UpdateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.dashboardService.UpdateTeamMember(instructorID, teamMemberID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team member updated successfully",
	})
}

// RemoveTeamMember handles DELETE /api/v1/instructor/team/:id
func (h *DashboardHandler) RemoveTeamMember(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	teamMemberIDStr := c.Param("id")
	teamMemberID, err := uuid.Parse(teamMemberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team member ID"})
		return
	}

	err = h.dashboardService.RemoveTeamMember(instructorID, teamMemberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove team member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team member removed successfully",
	})
}

// GetNotificationSettings handles GET /api/v1/instructor/notifications/settings
func (h *DashboardHandler) GetNotificationSettings(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	settings, err := h.dashboardService.GetNotificationSettings(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// UpdateNotificationSettings handles PUT /api/v1/instructor/notifications/settings
func (h *DashboardHandler) UpdateNotificationSettings(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var settings []model.NotificationSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.dashboardService.UpdateNotificationSettings(instructorID, settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification settings updated successfully",
	})
}

// GetCourse handles GET /api/v1/instructor/courses/{id}
func (h *DashboardHandler) GetCourse(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.dashboardService.GetCourse(instructorID, courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    course,
	})
}

// UpdateCourse handles PUT /api/v1/instructor/courses/{id}
func (h *DashboardHandler) UpdateCourse(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var req model.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	course, err := h.dashboardService.UpdateCourse(instructorID, courseID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    course,
	})
}

// DeleteCourse handles DELETE /api/v1/instructor/courses/{id}
func (h *DashboardHandler) DeleteCourse(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	err = h.dashboardService.DeleteCourse(instructorID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Course deleted successfully",
	})
}

// CreateCourse handles POST /api/v1/instructor/courses
func (h *DashboardHandler) CreateCourse(c *gin.Context) {
	instructorID, err := getInstructorIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	course, err := h.dashboardService.CreateCourse(instructorID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    course,
	})
}

// Helper functions


// parseIntParam parses integer parameter with default value
func parseIntParam(c *gin.Context, param string, defaultValue int) int {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
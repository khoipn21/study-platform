package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type InstructorDashboardHandler struct {
	serviceURL    string
	client        *http.Client
	courseHandler *CourseHandler
}

func NewInstructorDashboardHandler(serviceURL string, courseHandler *CourseHandler) *InstructorDashboardHandler {
	return &InstructorDashboardHandler{
		serviceURL:    serviceURL,
		client:        &http.Client{},
		courseHandler: courseHandler,
	}
}

// GetDashboardOverview handles GET /api/v1/instructor/dashboard/overview
func (h *InstructorDashboardHandler) GetDashboardOverview(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/dashboard/overview", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetInstructorCourses handles GET /api/v1/instructor/courses
func (h *InstructorDashboardHandler) GetInstructorCourses(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/courses", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// CreateCourse handles POST /api/v1/instructor/courses
func (h *InstructorDashboardHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("========== INSTRUCTOR DASHBOARD HANDLER CreateCourse CALLED - DELEGATING TO COURSE HANDLER ==========\n")

	// Delegate to the proper course handler that has lecture creation and Lemon Squeezy logic
	h.courseHandler.CreateCourse(w, r)
}

// GetRevenueAnalytics handles GET /api/v1/instructor/analytics/revenue
func (h *InstructorDashboardHandler) GetRevenueAnalytics(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/analytics/revenue", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetEngagementAnalytics handles GET /api/v1/instructor/analytics/engagement
func (h *InstructorDashboardHandler) GetEngagementAnalytics(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/analytics/engagement", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetStudents handles GET /api/v1/instructor/students
func (h *InstructorDashboardHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/students", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// BroadcastCommunication handles POST /api/v1/instructor/communication/broadcast
func (h *InstructorDashboardHandler) BroadcastCommunication(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("POST", "/api/v1/instructor/communication/broadcast", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetVideoAnalytics handles GET /api/v1/instructor/videos/analytics
func (h *InstructorDashboardHandler) GetVideoAnalytics(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/videos/analytics", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// BulkCourseOperations handles POST /api/v1/instructor/courses/{id}/bulk-operations
func (h *InstructorDashboardHandler) BulkCourseOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/api/v1/instructor/courses/%s/bulk-operations", courseID)
	result, err := h.proxyRequest("POST", path, r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetTeamMembers handles GET /api/v1/instructor/team
func (h *InstructorDashboardHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/team", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// InviteTeamMember handles POST /api/v1/instructor/team/invite
func (h *InstructorDashboardHandler) InviteTeamMember(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("POST", "/api/v1/instructor/team/invite", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// UpdateDashboardSettings handles PUT /api/v1/instructor/dashboard/settings
func (h *InstructorDashboardHandler) UpdateDashboardSettings(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("PUT", "/api/v1/instructor/dashboard/settings", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetCourseAnalytics handles GET /api/v1/instructor/courses/{id}/analytics
func (h *InstructorDashboardHandler) GetCourseAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/courses/%s/analytics", courseID)
	result, err := h.proxyRequest("GET", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetCourse handles GET /api/v1/instructor/courses/{id}
func (h *InstructorDashboardHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/courses/%s", courseID)
	result, err := h.proxyRequest("GET", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// UpdateCourse handles PUT /api/v1/instructor/courses/{id}
func (h *InstructorDashboardHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/api/v1/instructor/courses/%s", courseID)
	result, err := h.proxyRequest("PUT", path, r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// DeleteCourse handles DELETE /api/v1/instructor/courses/{id}
func (h *InstructorDashboardHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/courses/%s", courseID)
	result, err := h.proxyRequest("DELETE", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetStudentDetails handles GET /api/v1/instructor/students/{id}
func (h *InstructorDashboardHandler) GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/students/%s", studentID)
	result, err := h.proxyRequest("GET", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetStudentAnalytics handles GET /api/v1/instructor/analytics/students
func (h *InstructorDashboardHandler) GetStudentAnalytics(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/analytics/students", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetVideoEngagement handles GET /api/v1/instructor/videos/{id}/engagement
func (h *InstructorDashboardHandler) GetVideoEngagement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/videos/%s/engagement", videoID)
	result, err := h.proxyRequest("GET", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetCommunicationHistory handles GET /api/v1/instructor/communication/history
func (h *InstructorDashboardHandler) GetCommunicationHistory(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/communication/history", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// SetupAutomatedMessages handles POST /api/v1/instructor/communication/automated
func (h *InstructorDashboardHandler) SetupAutomatedMessages(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("POST", "/api/v1/instructor/communication/automated", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetAISuggestions handles GET /api/v1/instructor/suggestions
func (h *InstructorDashboardHandler) GetAISuggestions(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/suggestions", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// ImplementSuggestion handles POST /api/v1/instructor/suggestions/{id}/implement
func (h *InstructorDashboardHandler) ImplementSuggestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suggestionID := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/api/v1/instructor/suggestions/%s/implement", suggestionID)
	result, err := h.proxyRequest("POST", path, r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// UpdateTeamMember handles PUT /api/v1/instructor/team/{id}
func (h *InstructorDashboardHandler) UpdateTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	memberID := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/api/v1/instructor/team/%s", memberID)
	result, err := h.proxyRequest("PUT", path, r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// RemoveTeamMember handles DELETE /api/v1/instructor/team/{id}
func (h *InstructorDashboardHandler) RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	memberID := vars["id"]

	path := fmt.Sprintf("/api/v1/instructor/team/%s", memberID)
	result, err := h.proxyRequest("DELETE", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetNotifications handles GET /api/v1/instructor/notifications
func (h *InstructorDashboardHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/notifications", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetNotificationSettings handles GET /api/v1/instructor/notifications/settings
func (h *InstructorDashboardHandler) GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/api/v1/instructor/notifications/settings", r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// UpdateNotificationSettings handles PUT /api/v1/instructor/notifications/settings
func (h *InstructorDashboardHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("PUT", "/api/v1/instructor/notifications/settings", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// UploadVideo handles POST /api/v1/instructor/videos/upload
func (h *InstructorDashboardHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	result, err := h.proxyRequest("POST", "/api/v1/instructor/videos/upload", r, body)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// GetVideoStatus handles GET /api/v1/instructor/videos/status/{lectureId}
func (h *InstructorDashboardHandler) GetVideoStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lectureID := vars["lectureId"]

	path := fmt.Sprintf("/api/v1/instructor/videos/status/%s", lectureID)
	result, err := h.proxyRequest("GET", path, r, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	h.writeJSONResponse(w, result)
}

// Health check endpoint
func (h *InstructorDashboardHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyRequest("GET", "/health", r, nil)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	h.writeJSONResponse(w, result)
}

// Helper methods

func (h *InstructorDashboardHandler) proxyRequest(method, path string, originalReq *http.Request, body []byte) (*http.Response, error) {
	url := h.serviceURL + path

	// Add query parameters from original request
	if originalReq.URL.RawQuery != "" {
		url += "?" + originalReq.URL.RawQuery
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(originalReq.Context(), method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Copy headers from original request
	for key, values := range originalReq.Header {
		// Skip host header to avoid issues
		if strings.ToLower(key) == "host" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set content type for body requests
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return h.client.Do(req)
}

func (h *InstructorDashboardHandler) writeJSONResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}
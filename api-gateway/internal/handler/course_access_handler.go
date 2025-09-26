package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	coursepb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type CourseAccessHandler struct {
	courseClient coursepb.CourseServiceClient
	logger       logger.Logger
}

func NewCourseAccessHandler(courseConn *grpc.ClientConn, logger logger.Logger) *CourseAccessHandler {
	return &CourseAccessHandler{
		courseClient: coursepb.NewCourseServiceClient(courseConn),
		logger:       logger,
	}
}

// CheckCourseAccess checks if user has access to a specific course
func (h *CourseAccessHandler) CheckCourseAccess(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CheckCourseAccess request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get course ID from URL path
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		courseID = vars["course_id"] // Try alternative param name
	}
	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check enrollment
	enrollmentReq := &coursepb.GetEnrollmentRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	enrollment, err := h.courseClient.GetEnrollment(ctx, enrollmentReq)
	if err != nil {
		h.logger.Errorf("Failed to check enrollment: %v", err)
		// User not enrolled - return access denied
		response := map[string]interface{}{
			"success":    false,
			"hasAccess":  false,
			"message":    "Not enrolled in this course",
			"courseId":   courseID,
			"userId":     userID,
		}
		h.sendJSON(w, http.StatusForbidden, response)
		return
	}

	// Check if enrollment is active
	hasAccess := enrollment.Enrollment.Status == "active"
	statusCode := http.StatusOK
	if !hasAccess {
		statusCode = http.StatusForbidden
	}

	response := map[string]interface{}{
		"success":          true,
		"hasAccess":        hasAccess,
		"enrollmentStatus": enrollment.Enrollment.Status,
		"courseId":         courseID,
		"userId":           userID,
		"enrolledAt":       enrollment.Enrollment.EnrolledAt.AsTime().Format(time.RFC3339),
		"progressPercent":  enrollment.Enrollment.ProgressPercentage,
	}

	if enrollment.Enrollment.LastAccessed != nil {
		response["lastAccessed"] = enrollment.Enrollment.LastAccessed.AsTime().Format(time.RFC3339)
	}

	if enrollment.Enrollment.CompletedAt != nil {
		response["completedAt"] = enrollment.Enrollment.CompletedAt.AsTime().Format(time.RFC3339)
	}

	h.sendJSON(w, statusCode, response)
}

// GetMyEnrolledCourses returns list of courses user is enrolled in
func (h *CourseAccessHandler) GetMyEnrolledCourses(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetMyEnrolledCourses request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user's enrollments
	enrollmentsReq := &coursepb.ListEnrollmentsRequest{
		UserId:   userID,
		Status:   "active", // Only active enrollments
		Page:     1,
		PageSize: 100, // Get all enrollments for now
	}

	enrollmentsResp, err := h.courseClient.ListEnrollments(ctx, enrollmentsReq)
	if err != nil {
		h.logger.Errorf("Failed to list enrollments: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve enrolled courses")
		return
	}

	// Get course details for each enrollment
	var enrolledCourses []map[string]interface{}

	for _, enrollment := range enrollmentsResp.Enrollments {
		courseReq := &coursepb.GetCourseRequest{
			Id: enrollment.CourseId,
		}

		course, err := h.courseClient.GetCourse(ctx, courseReq)
		if err != nil {
			h.logger.Errorf("Failed to get course %s: %v", enrollment.CourseId, err)
			continue // Skip this course if we can't fetch it
		}

		courseData := map[string]interface{}{
			"id":                   course.Course.Id,
			"title":                course.Course.Title,
			"description":          course.Course.Description,
			"instructor_name":      course.Course.InstructorName,
			"category":             course.Course.Category,
			"difficulty_level":     course.Course.Level,
			"thumbnail_url":        course.Course.ThumbnailUrl,
			"duration_minutes":     course.Course.DurationMinutes,
			"rating":               course.Course.Rating,
			"rating_count":         course.Course.RatingCount,
			"enrollment_count":     course.Course.EnrollmentCount,
			"price":                course.Course.Price,
			"currency":             course.Course.Currency,
			"is_free":              course.Course.Price == 0,
			"access_type":          func() string { if course.Course.Price == 0 { return "free" } else { return "paid" } }(),
			"progress":             int(enrollment.ProgressPercentage),
			"certificate_available": true,
			"mobile_access":        true,
			"lifetime_access":      course.Course.Price > 0,
			"language":             "Tiếng Việt",
			"total_lectures":       len(course.Course.Lectures),
			"next_lecture_id":      getNextLectureId(course.Course.Lectures, enrollment.ProgressPercentage),
			"enrollment": map[string]interface{}{
				"id":          enrollment.Id,
				"status":      enrollment.Status,
				"enrolled_at": enrollment.EnrolledAt.AsTime().Format(time.RFC3339),
			},
		}

		if enrollment.LastAccessed != nil {
			courseData["last_accessed"] = enrollment.LastAccessed.AsTime().Format(time.RFC3339)
		}

		if enrollment.CompletedAt != nil {
			courseData["completed_at"] = enrollment.CompletedAt.AsTime().Format(time.RFC3339)
		}

		enrolledCourses = append(enrolledCourses, courseData)
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Enrolled courses retrieved successfully",
		"data": map[string]interface{}{
			"courses": enrolledCourses,
			"total":   len(enrolledCourses),
			"userId":  userID,
		},
	}

	h.sendJSON(w, http.StatusOK, response)
}

// GetEnrolledCourseLectures returns lectures for an enrolled course
func (h *CourseAccessHandler) GetEnrolledCourseLectures(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetEnrolledCourseLectures request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get course ID from URL path
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		courseID = vars["course_id"] // Try alternative param name
	}
	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check enrollment first
	enrollmentReq := &coursepb.GetEnrollmentRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	enrollment, err := h.courseClient.GetEnrollment(ctx, enrollmentReq)
	if err != nil {
		h.logger.Errorf("Failed to check enrollment: %v", err)
		h.sendError(w, http.StatusForbidden, "Access denied: Not enrolled in this course")
		return
	}

	// Check if enrollment is active
	if enrollment.Enrollment.Status != "active" {
		h.sendError(w, http.StatusForbidden, "Access denied: Enrollment is not active")
		return
	}

	// Get lectures for the course
	lecturesReq := &coursepb.ListLecturesRequest{
		CourseId: courseID,
		Status:   "published", // Only published lectures
		Page:     1,
		PageSize: 100, // Get all lectures for now
	}

	lecturesResp, err := h.courseClient.ListLectures(ctx, lecturesReq)
	if err != nil {
		h.logger.Errorf("Failed to list lectures: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve course lectures")
		return
	}

	// Convert lectures to response format
	var lectures []map[string]interface{}
	for _, lecture := range lecturesResp.Lectures {
		lectureData := map[string]interface{}{
			"lectureId":        lecture.Id,
			"title":            lecture.Title,
			"description":      lecture.Description,
			"type":             lecture.Type,
			"orderNumber":      lecture.OrderNumber,
			"durationMinutes":  lecture.DurationMinutes,
			"videoUrl":         lecture.VideoUrl,
			"videoId":          lecture.VideoId,
			"isFree":           lecture.IsFree,
			"status":           lecture.Status,
		}

		lectures = append(lectures, lectureData)
	}

	response := map[string]interface{}{
		"success":          true,
		"lectures":         lectures,
		"totalCount":       len(lectures),
		"courseId":         courseID,
		"userId":           userID,
		"enrollmentStatus": enrollment.Enrollment.Status,
	}

	h.sendJSON(w, http.StatusOK, response)
}

// Helper methods
func (h *CourseAccessHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Errorf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *CourseAccessHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"success": false,
		"error":   message,
		"message": message,
	}
	h.sendJSON(w, statusCode, response)
}

// getNextLectureId determines the next lecture based on progress percentage
func getNextLectureId(lectures []*coursepb.Lecture, progressPercentage float64) string {
	if len(lectures) == 0 {
		return ""
	}

	// If no progress, return first lecture
	if progressPercentage == 0 {
		return lectures[0].Id
	}

	// If 100% complete, return last lecture
	if progressPercentage >= 100 {
		return lectures[len(lectures)-1].Id
	}

	// Calculate which lecture to continue from based on progress
	lectureIndex := int((progressPercentage / 100.0) * float64(len(lectures)))

	// Ensure we don't go out of bounds
	if lectureIndex >= len(lectures) {
		lectureIndex = len(lectures) - 1
	}

	return lectures[lectureIndex].Id
}
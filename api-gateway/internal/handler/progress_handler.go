package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	progresspb "github.com/study-platform/progress-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
)

type ProgressHandler struct {
	progressClient progresspb.ProgressServiceClient
	logger         logger.Logger
}

func NewProgressHandler(progressConn *grpc.ClientConn, logger logger.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressClient: progresspb.NewProgressServiceClient(progressConn),
		logger:         logger,
	}
}

type UpdateProgressRequest struct {
	CourseID           string  `json:"course_id"`
	LectureID          string  `json:"lecture_id"`
	ProgressPercentage float64 `json:"progress_percentage"`
	WatchTimeSeconds   int32   `json:"watch_time_seconds"`
	IsCompleted        bool    `json:"is_completed"`
}

type CreateEnrollmentRequest struct {
	CourseID string `json:"course_id"`
}

type UpdateEnrollmentStatusRequest struct {
	CourseID string `json:"course_id"`
	Status   string `json:"status"`
}

type MarkLectureCompleteRequest struct {
	CourseID         string `json:"course_id"`
	LectureID        string `json:"lecture_id"`
	WatchTimeSeconds int32  `json:"watch_time_seconds"`
}

func (h *ProgressHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CourseID == "" || req.LectureID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID and lecture ID are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.UpdateProgressRequest{
		UserId:             userID,
		CourseId:           req.CourseID,
		LectureId:          req.LectureID,
		ProgressPercentage: req.ProgressPercentage,
		WatchTimeSeconds:   req.WatchTimeSeconds,
		IsCompleted:        req.IsCompleted,
	}

	resp, err := h.progressClient.UpdateProgress(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to update progress")
		return
	}

	// Format response
	data := h.formatUserProgress(resp.Progress)
	h.sendSuccess(w, resp.Message, data)
}

func (h *ProgressHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]
	lectureID := vars["lecture_id"]

	if courseID == "" || lectureID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID and lecture ID are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetProgressRequest{
		UserId:    userID,
		CourseId:  courseID,
		LectureId: lectureID,
	}

	resp, err := h.progressClient.GetProgress(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get progress")
		return
	}

	// Format response
	data := h.formatUserProgress(resp.Progress)
	h.sendSuccess(w, "Progress retrieved successfully", data)
}

func (h *ProgressHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetUserProgressRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	resp, err := h.progressClient.GetUserProgress(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get user progress")
		return
	}

	// Format response
	var progressList []map[string]interface{}
	for _, progress := range resp.Progress {
		progressList = append(progressList, h.formatUserProgress(progress))
	}

	data := map[string]interface{}{
		"progress":                    progressList,
		"overall_progress_percentage": resp.OverallProgressPercentage,
	}

	h.sendSuccess(w, "User progress retrieved successfully", data)
}

func (h *ProgressHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CourseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.CreateEnrollmentRequest{
		UserId:   userID,
		CourseId: req.CourseID,
	}

	resp, err := h.progressClient.CreateEnrollment(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to create enrollment")
		return
	}

	// Format response
	data := h.formatEnrollment(resp.Enrollment)
	h.sendSuccess(w, resp.Message, data)
}

func (h *ProgressHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetEnrollmentRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	resp, err := h.progressClient.GetEnrollment(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get enrollment")
		return
	}

	// Format response
	data := h.formatEnrollment(resp.Enrollment)
	h.sendSuccess(w, "Enrollment retrieved successfully", data)
}

func (h *ProgressHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	status := r.URL.Query().Get("status")

	// Set defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.ListEnrollmentsRequest{
		UserId:   userID,
		Status:   status,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := h.progressClient.ListEnrollments(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to list enrollments")
		return
	}

	// Format response
	var enrollments []map[string]interface{}
	for _, enrollment := range resp.Enrollments {
		enrollments = append(enrollments, h.formatEnrollment(enrollment))
	}

	data := map[string]interface{}{
		"enrollments": enrollments,
		"total":       resp.TotalCount,
		"page":        resp.Page,
		"page_size":   resp.PageSize,
		"total_pages": (resp.TotalCount + resp.PageSize - 1) / resp.PageSize,
	}

	h.sendSuccess(w, "Enrollments retrieved successfully", data)
}

func (h *ProgressHandler) MarkLectureComplete(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req MarkLectureCompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CourseID == "" || req.LectureID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID and lecture ID are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.MarkLectureCompleteRequest{
		UserId:           userID,
		CourseId:         req.CourseID,
		LectureId:        req.LectureID,
		WatchTimeSeconds: req.WatchTimeSeconds,
	}

	resp, err := h.progressClient.MarkLectureComplete(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to mark lecture complete")
		return
	}

	// Format response
	data := map[string]interface{}{
		"progress":         h.formatUserProgress(resp.Progress),
		"course_completed": resp.CourseCompleted,
	}

	h.sendSuccess(w, resp.Message, data)
}

func (h *ProgressHandler) GetLectureProgress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetLectureProgressRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	resp, err := h.progressClient.GetLectureProgress(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get lecture progress")
		return
	}

	// Format response
	var lectureProgress []map[string]interface{}
	for _, lp := range resp.LectureProgress {
		lectureProgress = append(lectureProgress, h.formatLectureProgress(lp))
	}

	data := map[string]interface{}{
		"lecture_progress": lectureProgress,
	}

	h.sendSuccess(w, "Lecture progress retrieved successfully", data)
}

func (h *ProgressHandler) GetCourseCompletion(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetCourseCompletionRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	resp, err := h.progressClient.GetCourseCompletion(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get course completion")
		return
	}

	// Format response
	data := h.formatCourseCompletion(resp.Completion)
	h.sendSuccess(w, "Course completion retrieved successfully", data)
}

func (h *ProgressHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call progress service
	grpcReq := &progresspb.GetUserAnalyticsRequest{
		UserId: userID,
	}

	resp, err := h.progressClient.GetUserAnalytics(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get user analytics")
		return
	}

	// Format response
	data := h.formatUserAnalytics(resp.Analytics)
	h.sendSuccess(w, "User analytics retrieved successfully", data)
}

// Helper methods
func (h *ProgressHandler) formatUserProgress(progress *progresspb.UserProgress) map[string]interface{} {
	data := map[string]interface{}{
		"id":                  progress.Id,
		"user_id":             progress.UserId,
		"course_id":           progress.CourseId,
		"lecture_id":          progress.LectureId,
		"progress_percentage": progress.ProgressPercentage,
		"watch_time_seconds":  progress.WatchTimeSeconds,
		"is_completed":        progress.IsCompleted,
		"last_accessed":       progress.LastAccessed.AsTime(),
		"created_at":          progress.CreatedAt.AsTime(),
		"updated_at":          progress.UpdatedAt.AsTime(),
	}

	if progress.CompletedAt != nil {
		data["completed_at"] = progress.CompletedAt.AsTime()
	}

	return data
}

func (h *ProgressHandler) formatEnrollment(enrollment *progresspb.Enrollment) map[string]interface{} {
	data := map[string]interface{}{
		"id":                     enrollment.Id,
		"user_id":                enrollment.UserId,
		"course_id":              enrollment.CourseId,
		"status":                 enrollment.Status,
		"progress_percentage":    enrollment.ProgressPercentage,
		"completed_lectures":     enrollment.CompletedLectures,
		"total_lectures":         enrollment.TotalLectures,
		"total_watch_time_seconds": enrollment.TotalWatchTimeSeconds,
		"enrolled_at":            enrollment.EnrolledAt.AsTime(),
		"created_at":             enrollment.CreatedAt.AsTime(),
		"updated_at":             enrollment.UpdatedAt.AsTime(),
	}

	if enrollment.CompletedAt != nil {
		data["completed_at"] = enrollment.CompletedAt.AsTime()
	}

	if enrollment.LastAccessed != nil {
		data["last_accessed"] = enrollment.LastAccessed.AsTime()
	}

	return data
}

func (h *ProgressHandler) formatLectureProgress(lp *progresspb.LectureProgress) map[string]interface{} {
	data := map[string]interface{}{
		"lecture_id":          lp.LectureId,
		"title":               lp.Title,
		"order_number":        lp.OrderNumber,
		"progress_percentage": lp.ProgressPercentage,
		"watch_time_seconds":  lp.WatchTimeSeconds,
		"is_completed":        lp.IsCompleted,
	}

	if lp.LastAccessed != nil {
		data["last_accessed"] = lp.LastAccessed.AsTime()
	}

	if lp.CompletedAt != nil {
		data["completed_at"] = lp.CompletedAt.AsTime()
	}

	return data
}

func (h *ProgressHandler) formatCourseCompletion(completion *progresspb.CourseCompletion) map[string]interface{} {
	var lectureProgress []map[string]interface{}
	for _, lp := range completion.LectureProgress {
		lectureProgress = append(lectureProgress, h.formatLectureProgress(lp))
	}

	data := map[string]interface{}{
		"course_id":               completion.CourseId,
		"course_title":            completion.CourseTitle,
		"user_id":                 completion.UserId,
		"completion_percentage":   completion.CompletionPercentage,
		"completed_lectures":      completion.CompletedLectures,
		"total_lectures":          completion.TotalLectures,
		"total_watch_time_seconds": completion.TotalWatchTimeSeconds,
		"lecture_progress":        lectureProgress,
	}

	if completion.StartedAt != nil {
		data["started_at"] = completion.StartedAt.AsTime()
	}

	if completion.CompletedAt != nil {
		data["completed_at"] = completion.CompletedAt.AsTime()
	}

	if completion.LastAccessed != nil {
		data["last_accessed"] = completion.LastAccessed.AsTime()
	}

	return data
}

func (h *ProgressHandler) formatUserAnalytics(analytics *progresspb.UserAnalytics) map[string]interface{} {
	data := map[string]interface{}{
		"user_id":                     analytics.UserId,
		"total_courses_enrolled":      analytics.TotalCoursesEnrolled,
		"total_courses_completed":     analytics.TotalCoursesCompleted,
		"total_lectures_completed":    analytics.TotalLecturesCompleted,
		"total_watch_time_seconds":    analytics.TotalWatchTimeSeconds,
		"average_progress_percentage": analytics.AverageProgressPercentage,
		"courses_in_progress":         analytics.CoursesInProgress,
		"most_active_day":             analytics.MostActiveDay,
		"streak_days":                 analytics.StreakDays,
	}

	if analytics.LastActivity != nil {
		data["last_activity"] = analytics.LastActivity.AsTime()
	}

	return data
}

func (h *ProgressHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *ProgressHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *ProgressHandler) handleGRPCError(w http.ResponseWriter, err error, defaultMessage string) {
	// Same implementation as auth handler
	h.logger.Errorf("gRPC error: %v", err)
	h.sendError(w, http.StatusInternalServerError, defaultMessage)
}
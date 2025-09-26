package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	progresspb "github.com/study-platform/progress-service/proto"
	"github.com/study-platform/pkg/logger"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type ProgressTrackingHandler struct {
	progressClient progresspb.ProgressServiceClient
	logger         logger.Logger
}

func NewProgressTrackingHandler(progressConn *grpc.ClientConn, logger logger.Logger) *ProgressTrackingHandler {
	return &ProgressTrackingHandler{
		progressClient: progresspb.NewProgressServiceClient(progressConn),
		logger:         logger,
	}
}

// TrackProgress tracks user progress on lectures
func (h *ProgressTrackingHandler) TrackProgress(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("TrackProgress request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request body
	var req struct {
		LectureID        string  `json:"lecture_id"`
		CourseID         string  `json:"course_id"`
		WatchTimeSeconds int32   `json:"watch_time_seconds"`
		ProgressPercent  float32 `json:"progress_percentage"`
		IsCompleted      bool    `json:"is_completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.LectureID == "" || req.CourseID == "" {
		h.sendError(w, http.StatusBadRequest, "lecture_id and course_id are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Track progress
	progressReq := &progresspb.UpdateProgressRequest{
		UserId:             userID,
		CourseId:           req.CourseID,
		LectureId:          req.LectureID,
		WatchTimeSeconds:   req.WatchTimeSeconds,
		ProgressPercentage: float64(req.ProgressPercent),
		IsCompleted:        req.IsCompleted,
	}

	progressResp, err := h.progressClient.UpdateProgress(ctx, progressReq)
	if err != nil {
		h.logger.Errorf("Failed to update progress: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to track progress")
		return
	}

	response := map[string]interface{}{
		"success":         true,
		"message":         "Progress tracked successfully",
		"userId":          userID,
		"courseId":        req.CourseID,
		"lectureId":       req.LectureID,
		"progressPercent": progressResp.Progress.ProgressPercentage,
		"isCompleted":     progressResp.Progress.IsCompleted,
		"updatedAt":       time.Now().Format(time.RFC3339),
	}

	h.sendJSON(w, http.StatusOK, response)
}

// MarkLectureComplete marks a lecture as completed
func (h *ProgressTrackingHandler) MarkLectureComplete(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("MarkLectureComplete request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request body
	var req struct {
		LectureID string `json:"lecture_id"`
		CourseID  string `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.LectureID == "" || req.CourseID == "" {
		h.sendError(w, http.StatusBadRequest, "lecture_id and course_id are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mark lecture as complete (100% progress)
	progressReq := &progresspb.UpdateProgressRequest{
		UserId:             userID,
		CourseId:           req.CourseID,
		LectureId:          req.LectureID,
		WatchTimeSeconds:   0, // Not updating watch time in this case
		ProgressPercentage: 100.0,
		IsCompleted:        true,
	}

	progressResp, err := h.progressClient.UpdateProgress(ctx, progressReq)
	if err != nil {
		h.logger.Errorf("Failed to mark lecture complete: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to mark lecture as complete")
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"message":     "Lecture marked as complete",
		"userId":      userID,
		"courseId":    req.CourseID,
		"lectureId":   req.LectureID,
		"isCompleted": progressResp.Progress.IsCompleted,
		"completedAt": time.Now().Format(time.RFC3339),
	}

	h.sendJSON(w, http.StatusOK, response)
}

// GetUserProgress gets user progress for a specific course
func (h *ProgressTrackingHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetUserProgress request received")

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

	// Get user progress
	progressReq := &progresspb.GetUserProgressRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	progressResp, err := h.progressClient.GetUserProgress(ctx, progressReq)
	if err != nil {
		h.logger.Errorf("Failed to get user progress: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve progress")
		return
	}

	// Convert progress entries to response format
	var lectureProgress []map[string]interface{}
	for _, progress := range progressResp.Progress {
		progressData := map[string]interface{}{
			"lectureId":         progress.LectureId,
			"watchTimeSeconds":  progress.WatchTimeSeconds,
			"progressPercent":   progress.ProgressPercentage,
			"isCompleted":       progress.IsCompleted,
			"updatedAt":         progress.UpdatedAt.AsTime().Format(time.RFC3339),
		}

		if progress.CompletedAt != nil {
			progressData["completedAt"] = progress.CompletedAt.AsTime().Format(time.RFC3339)
		}

		lectureProgress = append(lectureProgress, progressData)
	}

	response := map[string]interface{}{
		"success":         true,
		"userId":          userID,
		"courseId":        courseID,
		"lectureProgress": lectureProgress,
		"totalLectures":   len(lectureProgress),
	}

	// Calculate overall course progress
	if len(lectureProgress) > 0 {
		totalProgress := float64(0)
		completedLectures := 0

		for _, progress := range progressResp.Progress {
			totalProgress += progress.ProgressPercentage
			if progress.IsCompleted {
				completedLectures++
			}
		}

		avgProgress := totalProgress / float64(len(lectureProgress))
		response["overallProgress"] = avgProgress
		response["completedLectures"] = completedLectures
	}

	h.sendJSON(w, http.StatusOK, response)
}

// GetLectureProgress gets progress for a specific lecture
func (h *ProgressTrackingHandler) GetLectureProgress(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetLectureProgress request received")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get lecture and course IDs from URL path
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	if courseID == "" {
		courseID = vars["course_id"]
	}
	lectureID := vars["lectureId"]
	if lectureID == "" {
		lectureID = vars["lecture_id"]
	}

	if courseID == "" || lectureID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID and Lecture ID are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get lecture progress
	progressReq := &progresspb.GetProgressRequest{
		UserId:    userID,
		CourseId:  courseID,
		LectureId: lectureID,
	}

	progressResp, err := h.progressClient.GetProgress(ctx, progressReq)
	if err != nil {
		h.logger.Errorf("Failed to get lecture progress: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve lecture progress")
		return
	}

	response := map[string]interface{}{
		"success":          true,
		"userId":           userID,
		"courseId":         courseID,
		"lectureId":        lectureID,
		"watchTimeSeconds": progressResp.Progress.WatchTimeSeconds,
		"progressPercent":  progressResp.Progress.ProgressPercentage,
		"isCompleted":      progressResp.Progress.IsCompleted,
		"updatedAt":        progressResp.Progress.UpdatedAt.AsTime().Format(time.RFC3339),
	}

	if progressResp.Progress.CompletedAt != nil {
		response["completedAt"] = progressResp.Progress.CompletedAt.AsTime().Format(time.RFC3339)
	}

	h.sendJSON(w, http.StatusOK, response)
}

// Helper methods
func (h *ProgressTrackingHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Errorf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *ProgressTrackingHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"success": false,
		"error":   message,
		"message": message,
	}
	h.sendJSON(w, statusCode, response)
}
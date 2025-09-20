package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	coursepb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
)

type CourseHandler struct {
	courseClient coursepb.CourseServiceClient
	logger       logger.Logger
}

func NewCourseHandler(courseConn *grpc.ClientConn, logger logger.Logger) *CourseHandler {
	return &CourseHandler{
		courseClient: coursepb.NewCourseServiceClient(courseConn),
		logger:       logger,
	}
}

type CreateCourseRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	InstructorID string   `json:"instructor_id"`
	Category     string   `json:"category"`
	Level        string   `json:"level"`
	Price        float64  `json:"price"`
	Currency     string   `json:"currency"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags"`
}

type UpdateCourseRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Level        string   `json:"level"`
	Price        float64  `json:"price"`
	Currency     string   `json:"currency"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags"`
}

type CreateLectureRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	OrderNumber     int32  `json:"order_number"`
	DurationMinutes int32  `json:"duration_minutes"`
	VideoURL        string `json:"video_url"`
	VideoID         string `json:"video_id"`
	IsFree          bool   `json:"is_free"`
}

type UpdateLectureRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	OrderNumber     int32  `json:"order_number"`
	DurationMinutes int32  `json:"duration_minutes"`
	VideoURL        string `json:"video_url"`
	VideoID         string `json:"video_id"`
	Status          string `json:"status"`
	IsFree          bool   `json:"is_free"`
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.InstructorID == "" {
		h.sendError(w, http.StatusBadRequest, "Title, description, and instructor_id are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.CreateCourseRequest{
		Title:        req.Title,
		Description:  req.Description,
		InstructorId: req.InstructorID,
		Category:     req.Category,
		Level:        req.Level,
		Price:        req.Price,
		Currency:     req.Currency,
		ThumbnailUrl: req.ThumbnailURL,
		Status:       req.Status,
		Tags:         req.Tags,
	}

	resp, err := h.courseClient.CreateCourse(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to create course")
		return
	}

	// Format response
	data := h.formatCourse(resp.Course)
	h.sendSuccess(w, resp.Message, data)
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.GetCourseRequest{
		Id: courseID,
	}

	resp, err := h.courseClient.GetCourse(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get course")
		return
	}

	// Format response
	data := h.formatCourse(resp.Course)
	h.sendSuccess(w, "Course retrieved successfully", data)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	var req UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.UpdateCourseRequest{
		Id:           courseID,
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Level:        req.Level,
		Price:        req.Price,
		Currency:     req.Currency,
		ThumbnailUrl: req.ThumbnailURL,
		Status:       req.Status,
		Tags:         req.Tags,
	}

	resp, err := h.courseClient.UpdateCourse(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to update course")
		return
	}

	// Format response
	data := h.formatCourse(resp.Course)
	h.sendSuccess(w, resp.Message, data)
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.DeleteCourseRequest{
		Id: courseID,
	}

	resp, err := h.courseClient.DeleteCourse(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to delete course")
		return
	}

	h.sendSuccess(w, resp.Message, nil)
}

func (h *CourseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	category := r.URL.Query().Get("category")
	level := r.URL.Query().Get("level")
	status := r.URL.Query().Get("status")
	instructorID := r.URL.Query().Get("instructor_id")

	// Set defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.ListCoursesRequest{
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Category:     category,
		Level:        level,
		Status:       status,
		InstructorId: instructorID,
	}

	resp, err := h.courseClient.ListCourses(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to list courses")
		return
	}

	// Format response
	var courses []map[string]interface{}
	for _, course := range resp.Courses {
		courses = append(courses, h.formatCourse(course))
	}

	data := map[string]interface{}{
		"courses":    courses,
		"total":      resp.TotalCount,
		"page":       resp.Page,
		"page_size":  resp.PageSize,
		"total_pages": (resp.TotalCount + resp.PageSize - 1) / resp.PageSize,
	}

	h.sendSuccess(w, "Courses retrieved successfully", data)
}

func (h *CourseHandler) SearchCourses(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	category := r.URL.Query().Get("category")
	level := r.URL.Query().Get("level")
	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)
	minRating, _ := strconv.ParseFloat(r.URL.Query().Get("min_rating"), 64)

	// Set defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.SearchCoursesRequest{
		Query:     query,
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Category:  category,
		Level:     level,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
		MinRating: minRating,
	}

	resp, err := h.courseClient.SearchCourses(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to search courses")
		return
	}

	// Format response
	var courses []map[string]interface{}
	for _, course := range resp.Courses {
		courses = append(courses, h.formatCourse(course))
	}

	data := map[string]interface{}{
		"courses":    courses,
		"total":      resp.TotalCount,
		"page":       resp.Page,
		"page_size":  resp.PageSize,
		"total_pages": (resp.TotalCount + resp.PageSize - 1) / resp.PageSize,
	}

	h.sendSuccess(w, "Course search completed successfully", data)
}

func (h *CourseHandler) CreateLecture(w http.ResponseWriter, r *http.Request) {
	// Extract course ID from URL parameter
	vars := mux.Vars(r)
	courseID := vars["course_id"]
	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	var req CreateLectureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		h.sendError(w, http.StatusBadRequest, "Title is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.CreateLectureRequest{
		CourseId:        courseID,
		Title:           req.Title,
		Description:     req.Description,
		OrderNumber:     req.OrderNumber,
		DurationMinutes: req.DurationMinutes,
		VideoUrl:        req.VideoURL,
		VideoId:         req.VideoID,
		IsFree:          req.IsFree,
	}

	resp, err := h.courseClient.CreateLecture(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to create lecture")
		return
	}

	// Format response
	data := h.formatLecture(resp.Lecture)
	h.sendSuccess(w, resp.Message, data)
}

func (h *CourseHandler) GetLecture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lectureID := vars["id"]

	if lectureID == "" {
		h.sendError(w, http.StatusBadRequest, "Lecture ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.GetLectureRequest{
		Id: lectureID,
	}

	resp, err := h.courseClient.GetLecture(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get lecture")
		return
	}

	// Format response
	data := h.formatLecture(resp.Lecture)
	h.sendSuccess(w, "Lecture retrieved successfully", data)
}

func (h *CourseHandler) ListLectures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
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

	// Call course service
	grpcReq := &coursepb.ListLecturesRequest{
		CourseId: courseID,
		Status:   status,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := h.courseClient.ListLectures(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to list lectures")
		return
	}

	// Format response
	var lectures []map[string]interface{}
	for _, lecture := range resp.Lectures {
		lectures = append(lectures, h.formatLecture(lecture))
	}

	data := map[string]interface{}{
		"lectures":   lectures,
		"total":      resp.TotalCount,
		"page":       resp.Page,
		"page_size":  resp.PageSize,
		"total_pages": (resp.TotalCount + resp.PageSize - 1) / resp.PageSize,
	}

	h.sendSuccess(w, "Lectures retrieved successfully", data)
}

func (h *CourseHandler) EnrollInCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.EnrollInCourseRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	resp, err := h.courseClient.EnrollInCourse(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to enroll in course")
		return
	}

	// Format response
	data := h.formatEnrollment(resp.Enrollment)
	h.sendSuccess(w, resp.Message, data)
}

// Helper methods
func (h *CourseHandler) formatCourse(course *coursepb.Course) map[string]interface{} {
	return map[string]interface{}{
		"id":                course.Id,
		"title":             course.Title,
		"description":       course.Description,
		"instructor_id":     course.InstructorId,
		"instructor_name":   course.InstructorName,
		"category":          course.Category,
		"level":             course.Level,
		"price":             course.Price,
		"currency":          course.Currency,
		"thumbnail_url":     course.ThumbnailUrl,
		"status":            course.Status,
		"duration_minutes":  course.DurationMinutes,
		"enrollment_count":  course.EnrollmentCount,
		"rating":            course.Rating,
		"rating_count":      course.RatingCount,
		"tags":              course.Tags,
		"created_at":        course.CreatedAt.AsTime(),
		"updated_at":        course.UpdatedAt.AsTime(),
	}
}

func (h *CourseHandler) formatLecture(lecture *coursepb.Lecture) map[string]interface{} {
	return map[string]interface{}{
		"id":               lecture.Id,
		"course_id":        lecture.CourseId,
		"title":            lecture.Title,
		"description":      lecture.Description,
		"order_number":     lecture.OrderNumber,
		"duration_minutes": lecture.DurationMinutes,
		"video_url":        lecture.VideoUrl,
		"video_id":         lecture.VideoId,
		"status":           lecture.Status,
		"is_free":          lecture.IsFree,
		"created_at":       lecture.CreatedAt.AsTime(),
		"updated_at":       lecture.UpdatedAt.AsTime(),
	}
}

func (h *CourseHandler) formatEnrollment(enrollment *coursepb.Enrollment) map[string]interface{} {
	data := map[string]interface{}{
		"id":                  enrollment.Id,
		"user_id":             enrollment.UserId,
		"course_id":           enrollment.CourseId,
		"status":              enrollment.Status,
		"progress_percentage": enrollment.ProgressPercentage,
		"enrolled_at":         enrollment.EnrolledAt.AsTime(),
		"last_accessed":       enrollment.LastAccessed.AsTime(),
	}

	if enrollment.CompletedAt != nil {
		data["completed_at"] = enrollment.CompletedAt.AsTime()
	}

	return data
}

func (h *CourseHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *CourseHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *CourseHandler) handleGRPCError(w http.ResponseWriter, err error, defaultMessage string) {
	// Same implementation as auth handler
	h.logger.Errorf("gRPC error: %v", err)
	h.sendError(w, http.StatusInternalServerError, defaultMessage)
}
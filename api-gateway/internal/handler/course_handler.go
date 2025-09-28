package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	coursepb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
)

type CourseHandler struct {
	courseClient     coursepb.CourseServiceClient
	bucketServiceURL string
	logger           logger.Logger
}

func NewCourseHandler(courseConn *grpc.ClientConn, bucketServiceURL string, logger logger.Logger) *CourseHandler {
	return &CourseHandler{
		courseClient:     coursepb.NewCourseServiceClient(courseConn),
		bucketServiceURL: bucketServiceURL,
		logger:           logger,
	}
}

type CreateCourseRequest struct {
	Title                  string                     `json:"title"`
	Description            string                     `json:"description"`
	InstructorID           string                     `json:"instructor_id"`
	Category               string                     `json:"category"`
	Level                  string                     `json:"level"`
	Price                  float64                    `json:"price"`
	Currency               string                     `json:"currency"`
	ThumbnailURL           string                     `json:"thumbnail_url"`
	Status                 string                     `json:"status"`
	Tags                   []string                   `json:"tags"`
	Language               string                     `json:"language"`
	LearningOutcomes       []string                   `json:"learning_outcomes"`
	Requirements           []string                   `json:"requirements"`
	EstimatedDurationHours float64                    `json:"estimated_duration_hours"`
	AutoApproveEnrollment  bool                       `json:"auto_approve_enrollment"`
	AllowPreviews          bool                       `json:"allow_previews"`
	HasCertificate         bool                       `json:"has_certificate"`
	MobileAccess           bool                       `json:"mobile_access"`
	StartDate              string                     `json:"start_date"`
	EndDate                string                     `json:"end_date"`
	MaxStudents            int32                      `json:"max_students"`
	Lectures               []CreateLectureRequest     `json:"lectures"`
	Resources              []CreateResourceRequest    `json:"resources"`
}

type CreateResourceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	FileSize    int64  `json:"file_size"`
	IsPublic    bool   `json:"is_public"`
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
	ID              string                         `json:"id"`
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	Type            string                         `json:"type"`
	OrderNumber     int32                          `json:"order_number"`
	DurationMinutes int32                          `json:"duration_minutes"`
	VideoURL        string                         `json:"video_url"`
	VideoID         string                         `json:"video_id"`
	IsFree          bool                           `json:"is_free"`
	Resources       []CreateLectureResourceRequest `json:"resources"`
}

type CreateLectureResourceRequest struct {
	FileID       string `json:"file_id,omitempty"`
	Filename     string `json:"filename,omitempty"`
	OriginalName string `json:"original_name,omitempty"`
	FileType     string `json:"file_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	IsRequired   bool   `json:"is_required,omitempty"`
	IsPublic     bool   `json:"is_public,omitempty"`
	DownloadURL  string `json:"download_url,omitempty"` // Add this field to capture the URL
	UploadedAt   string `json:"uploaded_at,omitempty"`  // Add this field as well
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
	h.logger.Infof("========== COURSE HANDLER CreateCourse CALLED ==========")
	fmt.Printf("========== COURSE HANDLER CreateCourse CALLED ==========\n")

	// Read the raw body first for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		h.sendError(w, http.StatusBadRequest, "Could not read request body")
		return
	}

	h.logger.Infof("DEBUG: Raw request body: %s", string(bodyBytes))

	var req CreateCourseRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		h.logger.Errorf("Failed to decode request body: %v", err)
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	h.logger.Infof("Decoded request - Title: %s, Lectures count: %d", req.Title, len(req.Lectures))

	// TEMP DEBUG: Try to see if the lectures array is populated correctly
	lecturesJson, _ := json.Marshal(req.Lectures)
	h.logger.Infof("TEMP DEBUG: Lectures JSON: %s", string(lecturesJson))

	// Debug: Write lecture information to a debug file
	if len(req.Lectures) > 0 {
		h.logger.Infof("First lecture title: %s", req.Lectures[0].Title)
		h.logger.Infof("CRITICAL PARSE DEBUG: First lecture has %d resources", len(req.Lectures[0].Resources))
		fmt.Printf("CRITICAL PARSE DEBUG: First lecture has %d resources\n", len(req.Lectures[0].Resources))
		if len(req.Lectures[0].Resources) > 0 {
			for i, res := range req.Lectures[0].Resources {
				h.logger.Infof("PARSE DEBUG: Resource %d - Filename: %s, FileID: %s, DownloadURL: %s", i+1, res.Filename, res.FileID, res.DownloadURL)
				fmt.Printf("PARSE DEBUG: Resource %d - Filename: %s, FileID: %s, DownloadURL: %s\n", i+1, res.Filename, res.FileID, res.DownloadURL)
			}
		}
	} else {
		h.logger.Infof("NO LECTURES FOUND IN REQUEST")
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.InstructorID == "" {
		h.sendError(w, http.StatusBadRequest, "Title, description, and instructor_id are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	courseID := resp.Course.Id
	h.logger.Infof("Course created successfully: %s", courseID)

	// Debug: Log the number of lectures received
	h.logger.Infof("Number of lectures to create: %d", len(req.Lectures))

	// Create lectures for the course
	var createdLectures []*coursepb.Lecture
	for i, lecture := range req.Lectures {
		h.logger.Infof("Creating lecture %d: %s", i+1, lecture.Title)
		h.logger.Infof("CRITICAL DEBUG: Lecture %d has %d resources", i+1, len(lecture.Resources))
		fmt.Printf("CRITICAL DEBUG: Lecture %d has %d resources\n", i+1, len(lecture.Resources))
		if len(lecture.Resources) > 0 {
			for j, res := range lecture.Resources {
				h.logger.Infof("DEBUG: Resource %d - Filename: %s, FileID: %s, DownloadURL: %s", j+1, res.Filename, res.FileID, res.DownloadURL)
				fmt.Printf("DEBUG: Resource %d - Filename: %s, FileID: %s, DownloadURL: %s\n", j+1, res.Filename, res.FileID, res.DownloadURL)
			}
		} else {
			h.logger.Infof("CRITICAL: NO RESOURCES FOUND FOR LECTURE %s", lecture.Title)
			fmt.Printf("CRITICAL: NO RESOURCES FOUND FOR LECTURE %s\n", lecture.Title)
		}
		lectureReq := &coursepb.CreateLectureRequest{
			CourseId:        courseID,
			Title:           lecture.Title,
			Description:     lecture.Description,
			OrderNumber:     int32(i + 1), // Use sequential order
			DurationMinutes: lecture.DurationMinutes,
			VideoUrl:        lecture.VideoURL,
			VideoId:         lecture.VideoID,
			IsFree:          lecture.IsFree,
		}

		lectureResp, err := h.courseClient.CreateLecture(ctx, lectureReq)
		if err != nil {
			h.logger.Errorf("Failed to create lecture '%s': %v", lecture.Title, err)
			// Continue with other lectures, don't fail the entire course creation
			continue
		}
		createdLectures = append(createdLectures, lectureResp.Lecture)
		h.logger.Infof("Lecture created: %s (%s)", lectureResp.Lecture.Title, lectureResp.Lecture.Id)

		// Create resources for this lecture if any exist
		h.logger.Infof("DEBUG: Checking resources for lecture '%s' - resources count: %d", lecture.Title, len(lecture.Resources))
		if len(lecture.Resources) > 0 {
			h.logger.Infof("Creating %d resources for lecture: %s", len(lecture.Resources), lecture.Title)
			err = h.createLectureResources(ctx, lectureResp.Lecture.Id, lecture.Resources)
			if err != nil {
				h.logger.Errorf("Failed to create resources for lecture '%s': %v", lecture.Title, err)
				// Don't fail the entire course creation, just log the error
			}
		} else {
			h.logger.Infof("DEBUG: No resources found for lecture '%s'", lecture.Title)
		}
	}

	// Create Lemon Squeezy product if course has a price
	var lemonSqueezyProductID string
	if req.Price > 0 && req.Status == "published" {
		// Extract user information from context for authentication headers
		userID, userIDOk := r.Context().Value("user_id").(string)
		userRole := ""
		if userRoleValue := r.Context().Value("user_role"); userRoleValue != nil {
			userRole = userRoleValue.(string)
		}

		if !userIDOk || userID == "" {
			h.logger.Errorf("User ID not found in context for Lemon Squeezy product creation")
		} else {
			lemonSqueezyProductID, err = h.createLemonSqueezyProduct(ctx, resp.Course, userID, userRole)
			if err != nil {
				h.logger.Errorf("Failed to create Lemon Squeezy product for course %s: %v", courseID, err)
				// Don't fail the course creation, just log the error
			}
		}
	}

	// Format response with created lectures
	data := h.formatCourse(resp.Course)
	data["lectures"] = h.formatLectures(createdLectures)
	data["lemon_squeezy_product_id"] = lemonSqueezyProductID

	// Debug: Add debug information to response
	data["debug_info"] = map[string]interface{}{
		"lectures_requested": len(req.Lectures),
		"lectures_created": len(createdLectures),
	}

	h.sendSuccess(w, "MODIFIED API GATEWAY: Course and lectures created successfully", data)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing course first to check for changes
	getCourseReq := &coursepb.GetCourseRequest{Id: courseID}
	existingCourse, err := h.courseClient.GetCourse(ctx, getCourseReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get existing course")
		return
	}

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

	// Log if thumbnail changed (for S3 cleanup tracking)
	if existingCourse.Course.ThumbnailUrl != req.ThumbnailURL {
		h.logger.Infof("Course %s thumbnail changed: %s -> %s",
			courseID, existingCourse.Course.ThumbnailUrl, req.ThumbnailURL)
	}

	// Format response
	data := h.formatCourse(resp.Course)
	h.sendSuccess(w, resp.Message, data)
}

// UpdateCourseWithThumbnail handles course update with thumbnail file upload
func (h *CourseHandler) UpdateCourseWithThumbnail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	if courseID == "" {
		h.sendError(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Parse multipart form with a 32MB limit
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.logger.Errorf("Failed to parse multipart form: %v", err)
		h.sendError(w, http.StatusBadRequest, "Invalid multipart form")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get existing course to preserve unchanged fields
	getCourseReq := &coursepb.GetCourseRequest{Id: courseID}
	existingCourseResp, err := h.courseClient.GetCourse(ctx, getCourseReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get existing course")
		return
	}
	existingCourse := existingCourseResp.Course

	// Extract course data from form (use existing values as defaults)
	req := UpdateCourseRequest{
		Title:        getFormValueOrDefault(r, "title", existingCourse.Title),
		Description:  getFormValueOrDefault(r, "description", existingCourse.Description),
		Category:     getFormValueOrDefault(r, "category", existingCourse.Category),
		Level:        getFormValueOrDefault(r, "level", existingCourse.Level),
		Currency:     getFormValueOrDefault(r, "currency", existingCourse.Currency),
		Status:       getFormValueOrDefault(r, "status", existingCourse.Status),
		ThumbnailURL: existingCourse.ThumbnailUrl, // Will be updated if new file uploaded
	}

	// Parse optional price field
	if priceStr := r.FormValue("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			req.Price = price
		}
	} else {
		req.Price = existingCourse.Price
	}

	// Parse tags if provided
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
		// Trim whitespace from tags
		for i, tag := range req.Tags {
			req.Tags[i] = strings.TrimSpace(tag)
		}
	} else {
		req.Tags = existingCourse.Tags
	}

	// Handle thumbnail upload if provided
	var newThumbnailURL string
	oldThumbnailURL := existingCourse.ThumbnailUrl

	if file, fileHeader, err := r.FormFile("thumbnail"); err == nil {
		defer file.Close()

		h.logger.Infof("New thumbnail file uploaded for course %s: %s, size: %d bytes",
			courseID, fileHeader.Filename, fileHeader.Size)

		// Validate thumbnail file
		if err := h.validateThumbnailFile(fileHeader); err != nil {
			h.sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid thumbnail: %s", err.Error()))
			return
		}

		// Upload new thumbnail to bucket service
		uploadedURL, err := h.uploadThumbnailToBucket(ctx, file, fileHeader, r.Header.Get("Authorization"))
		if err != nil {
			h.logger.Errorf("Failed to upload thumbnail: %v", err)
			h.sendError(w, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}

		newThumbnailURL = uploadedURL
		req.ThumbnailURL = newThumbnailURL
		h.logger.Infof("New thumbnail uploaded successfully for course %s: %s", courseID, newThumbnailURL)
	}

	// Call course service to update
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

	// Log thumbnail change for S3 cleanup tracking
	if oldThumbnailURL != req.ThumbnailURL {
		h.logger.Infof("Course %s thumbnail updated: %s -> %s", courseID, oldThumbnailURL, req.ThumbnailURL)
		// TODO: Here you could call bucket service to clean up old thumbnail
		// This should be done asynchronously to avoid blocking the response
	}

	// Format response
	data := h.formatCourse(resp.Course)
	data["thumbnail_updated"] = newThumbnailURL != ""
	data["old_thumbnail_url"] = oldThumbnailURL

	h.sendSuccess(w, "Course updated successfully", data)
}

// Helper function to get form value or return default
func getFormValueOrDefault(r *http.Request, key, defaultValue string) string {
	if value := r.FormValue(key); value != "" {
		return value
	}
	return defaultValue
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

	h.logger.Infof("DEBUG: API Gateway ListLectures called for course_id=%s, page=%d, page_size=%d", courseID, page, pageSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call course service
	grpcReq := &coursepb.ListLecturesRequest{
		CourseId: courseID,
		Status:   status,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	h.logger.Infof("DEBUG: Calling course service gRPC ListLectures")
	resp, err := h.courseClient.ListLectures(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to list lectures")
		return
	}

	h.logger.Infof("DEBUG: Received gRPC response with %d lectures", len(resp.Lectures))

	// Format response
	var lectures []map[string]interface{}
	for i, lecture := range resp.Lectures {
		h.logger.Infof("DEBUG: Processing lecture %d (%s) with %d resources", i, lecture.Id, len(lecture.Resources))
		if len(lecture.Resources) > 0 {
			for j, res := range lecture.Resources {
				h.logger.Infof("DEBUG: Resource %d - ID: %s, Filename: %s, Type: %s", j, res.Id, res.Filename, res.FileType)
			}
		} else {
			h.logger.Infof("DEBUG: No resources found for lecture %s", lecture.Id)
		}
		formattedLecture := h.formatLecture(lecture)
		lectures = append(lectures, formattedLecture)
		h.logger.Infof("DEBUG: Formatted lecture %d has %d resources", i, len(formattedLecture["resources"].([]map[string]interface{})))
	}

	data := map[string]interface{}{
		"lectures":   lectures,
		"total":      resp.TotalCount,
		"page":       resp.Page,
		"page_size":  resp.PageSize,
		"total_pages": (resp.TotalCount + resp.PageSize - 1) / resp.PageSize,
	}

	h.logger.Infof("DEBUG: Final response data contains %d lectures", len(lectures))
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
	// Format resources
	var resources []map[string]interface{}
	for _, resource := range lecture.Resources {
		resources = append(resources, h.formatLectureResource(resource))
	}

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
		"resources":        resources,
	}
}

func (h *CourseHandler) formatLectureResource(resource *coursepb.LectureResource) map[string]interface{} {
	return map[string]interface{}{
		"id":            resource.Id,
		"lecture_id":    resource.LectureId,
		"file_id":       resource.FileId,
		"resource_type": resource.ResourceType,
		"display_order": resource.DisplayOrder,
		"is_required":   resource.IsRequired,
		"filename":      resource.Filename,
		"original_name": resource.OriginalName,
		"file_type":     resource.FileType,
		"file_size":     resource.FileSize,
		"download_url":  resource.DownloadUrl,
		"is_public":     resource.IsPublic,
		"created_at":    resource.CreatedAt.AsTime(),
		"updated_at":    resource.UpdatedAt.AsTime(),
		"uploaded_at":   resource.UploadedAt.AsTime(),
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

// formatLectures formats multiple lectures
func (h *CourseHandler) formatLectures(lectures []*coursepb.Lecture) []map[string]interface{} {
	var result []map[string]interface{}
	for _, lecture := range lectures {
		result = append(result, h.formatLecture(lecture))
	}
	return result
}

// createLemonSqueezyProduct creates a product in Lemon Squeezy for the course
func (h *CourseHandler) createLemonSqueezyProduct(ctx context.Context, course *coursepb.Course, userID, userRole string) (string, error) {
	h.logger.Infof("Creating Lemon Squeezy product for course: %s (user: %s, role: %s)", course.Title, userID, userRole)

	// Prepare request payload for payment service
	payload := map[string]interface{}{
		"course_id":    course.Id,
		"title":        course.Title,
		"description":  course.Description,
		"price":        course.Price,
		"currency":     course.Currency,
		"category":     course.Category,
		"instructor_id": course.InstructorId,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Call payment service to create Lemon Squeezy product
	paymentServiceURL := "http://payment-service:8088/api/v1/lemonsqueezy/products"
	req, err := http.NewRequestWithContext(ctx, "POST", paymentServiceURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	if userRole != "" {
		req.Header.Set("X-User-Role", userRole)
	}
	h.logger.Infof("Calling payment service with headers: X-User-ID=%s, X-User-Role=%s", userID, userRole)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call payment service: %v", err)
		return "", fmt.Errorf("failed to call payment service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		h.logger.Errorf("Payment service returned error: %d", resp.StatusCode)
		return "", fmt.Errorf("payment service error: status %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract product ID from response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	productID, ok := data["product_id"].(string)
	if !ok {
		return "", fmt.Errorf("product_id not found in response")
	}

	h.logger.Infof("Lemon Squeezy product created successfully: %s", productID)
	return productID, nil
}

// CreateCourseWithThumbnail handles course creation with thumbnail file upload
func (h *CourseHandler) CreateCourseWithThumbnail(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("========== COURSE HANDLER CreateCourseWithThumbnail CALLED ==========")

	// Parse multipart form with a 32MB limit
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.logger.Errorf("Failed to parse multipart form: %v", err)
		h.sendError(w, http.StatusBadRequest, "Invalid multipart form")
		return
	}

	// Extract course data from form
	req := CreateCourseRequest{
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
		InstructorID: r.FormValue("instructor_id"),
		Category:     r.FormValue("category"),
		Level:        r.FormValue("level"),
		Currency:     r.FormValue("currency"),
		Status:       r.FormValue("status"),
	}

	// Parse optional fields
	if priceStr := r.FormValue("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			req.Price = price
		}
	}

	// Parse tags if provided
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
		// Trim whitespace from tags
		for i, tag := range req.Tags {
			req.Tags[i] = strings.TrimSpace(tag)
		}
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.InstructorID == "" {
		h.sendError(w, http.StatusBadRequest, "Title, description, and instructor_id are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Handle thumbnail upload if provided
	var thumbnailURL string
	if file, fileHeader, err := r.FormFile("thumbnail"); err == nil {
		defer file.Close()

		h.logger.Infof("Thumbnail file uploaded: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

		// Validate thumbnail file
		if err := h.validateThumbnailFile(fileHeader); err != nil {
			h.sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid thumbnail: %s", err.Error()))
			return
		}

		// Upload thumbnail to bucket service
		uploadedURL, err := h.uploadThumbnailToBucket(ctx, file, fileHeader, r.Header.Get("Authorization"))
		if err != nil {
			h.logger.Errorf("Failed to upload thumbnail: %v", err)
			h.sendError(w, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}

		thumbnailURL = uploadedURL
		h.logger.Infof("Thumbnail uploaded successfully: %s", thumbnailURL)
	}

	// Set thumbnail URL in request
	req.ThumbnailURL = thumbnailURL

	// Call course service to create course
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

	courseID := resp.Course.Id
	h.logger.Infof("Course created successfully: %s", courseID)

	// Create Lemon Squeezy product if course has a price
	var lemonSqueezyProductID string
	if req.Price > 0 && req.Status == "published" {
		// Extract user information from context for authentication headers
		userID, userIDOk := r.Context().Value("user_id").(string)
		userRole := ""
		if userRoleValue := r.Context().Value("user_role"); userRoleValue != nil {
			userRole = userRoleValue.(string)
		}

		if !userIDOk || userID == "" {
			h.logger.Errorf("User ID not found in context for Lemon Squeezy product creation")
		} else {
			lemonSqueezyProductID, err = h.createLemonSqueezyProduct(ctx, resp.Course, userID, userRole)
			if err != nil {
				h.logger.Errorf("Failed to create Lemon Squeezy product for course %s: %v", courseID, err)
				// Don't fail the course creation, just log the error
			}
		}
	}

	// Format response
	data := h.formatCourse(resp.Course)
	data["lemon_squeezy_product_id"] = lemonSqueezyProductID
	data["thumbnail_uploaded"] = thumbnailURL != ""

	h.sendSuccess(w, "Course created successfully with thumbnail", data)
}

// uploadThumbnailToBucket uploads the thumbnail file to the bucket service
func (h *CourseHandler) uploadThumbnailToBucket(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, authHeader string) (string, error) {
	// Create a new HTTP request for bucket service
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Copy the file to the new form
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Reset file position to beginning
	file.Seek(0, 0)
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Add additional form fields
	writer.WriteField("bucket", "images") // Use images bucket for thumbnails
	writer.WriteField("is_public", "true") // Make thumbnails publicly accessible

	writer.Close()

	// Create HTTP request to bucket service
	bucketURL := fmt.Sprintf("http://%s/api/files/upload", h.bucketServiceURL)
	req, err := http.NewRequestWithContext(ctx, "POST", bucketURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to bucket service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bucket service returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response to get the file URL
	var uploadResponse struct {
		FileID       string `json:"file_id"`
		Filename     string `json:"filename"`
		URL          string `json:"url"`
		ThumbnailURL string `json:"thumbnail_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return "", fmt.Errorf("failed to decode bucket service response: %w", err)
	}

	// Return the thumbnail URL if available, otherwise the main URL
	if uploadResponse.ThumbnailURL != "" {
		return uploadResponse.ThumbnailURL, nil
	}
	return uploadResponse.URL, nil
}

// validateThumbnailFile validates the uploaded thumbnail file
func (h *CourseHandler) validateThumbnailFile(fileHeader *multipart.FileHeader) error {
	// Check file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		return fmt.Errorf("thumbnail file too large (max 5MB)")
	}

	// Check file extension
	filename := strings.ToLower(fileHeader.Filename)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	validExt := false
	for _, ext := range allowedExts {
		if strings.HasSuffix(filename, ext) {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("invalid file type (allowed: jpg, jpeg, png, gif, webp)")
	}

	return nil
}

// createLectureResources creates resources for a lecture by calling the course service HTTP API
func (h *CourseHandler) createLectureResources(ctx context.Context, lectureID string, resources []CreateLectureResourceRequest) error {
	h.logger.Infof("Creating %d resources for lecture: %s", len(resources), lectureID)

	for i, resource := range resources {
		// Extract file ID from download URL if not provided
		fileID := resource.FileID
		if fileID == "" && resource.DownloadURL != "" {
			fileID = h.extractFileIDFromURL(resource.DownloadURL)
		}

		// Infer resource type if not provided
		resourceType := resource.ResourceType
		if resourceType == "" {
			resourceType = h.inferResourceTypeFromFilename(resource.Filename)
		}

		h.logger.Infof("Creating resource %d: %s (file_id: %s, type: %s)", i+1, resource.Filename, fileID, resourceType)

		// Validate that we have a file ID
		if fileID == "" {
			h.logger.Errorf("Cannot create resource %s: missing file_id", resource.Filename)
			continue
		}

		// Create the resource payload
		resourcePayload := map[string]interface{}{
			"file_id":       fileID,
			"resource_type": resourceType,
			"display_order": i + 1, // Use sequential order
			"is_required":   resource.IsRequired,
		}

		payloadBytes, err := json.Marshal(resourcePayload)
		if err != nil {
			h.logger.Errorf("Failed to marshal resource payload: %v", err)
			continue // Skip this resource but continue with others
		}

		// Call the course service HTTP API to create the lecture resource
		courseServiceURL := "http://course-service:8092"
		apiURL := fmt.Sprintf("%s/lectures/%s/resources", courseServiceURL, lectureID)

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			h.logger.Errorf("Failed to create HTTP request: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			h.logger.Errorf("Failed to create lecture resource: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			h.logger.Infof("Successfully created resource for lecture %s: %s", lectureID, resource.Filename)
		} else {
			responseBody, _ := io.ReadAll(resp.Body)
			h.logger.Errorf("Failed to create lecture resource, status: %d, response: %s", resp.StatusCode, string(responseBody))
		}
	}

	return nil
}

// extractFileIDFromURL extracts the file ID from a download URL
// Expected URL format: https://bucket.com/path/to/file-id.extension
func (h *CourseHandler) extractFileIDFromURL(downloadURL string) string {
	if downloadURL == "" {
		return ""
	}

	// Extract filename from URL
	parts := strings.Split(downloadURL, "/")
	if len(parts) == 0 {
		return ""
	}

	filename := parts[len(parts)-1]

	// Remove extension to get file ID
	fileID := strings.TrimSuffix(filename, filepath.Ext(filename))

	h.logger.Infof("Extracted file ID '%s' from URL '%s'", fileID, downloadURL)
	return fileID
}

// inferResourceTypeFromFilename infers resource type from file extension
func (h *CourseHandler) inferResourceTypeFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "document"
	case ".ppt", ".pptx":
		return "slides"
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".mkv":
		return "video"
	case ".mp3", ".wav", ".aac":
		return "audio"
	case ".zip", ".rar", ".tar", ".gz":
		return "archive"
	case ".js", ".py", ".go", ".java", ".cpp", ".c", ".html", ".css":
		return "code"
	default:
		return "document" // Default fallback
	}
}

// GetLectureResourceDownloadURL proxies the request to course service for generating signed download URLs
func (h *CourseHandler) GetLectureResourceDownloadURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Prepare request to course service
	courseServiceURL := "http://course-service:8092"
	apiURL := fmt.Sprintf("%s/lecture-resources/%s/download-url", courseServiceURL, resourceID)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request to course service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add user ID to context for course service
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Content-Type", "application/json")

	// Forward the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Failed to generate download URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// GetLectureResourcePreviewURL proxies the request to course service for generating signed preview URLs
func (h *CourseHandler) GetLectureResourcePreviewURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Prepare request to course service
	courseServiceURL := "http://course-service:8092"
	apiURL := fmt.Sprintf("%s/lecture-resources/%s/preview-url", courseServiceURL, resourceID)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request to course service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add user ID to context for course service
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Content-Type", "application/json")

	// Forward the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Failed to generate preview URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}
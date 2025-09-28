package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/service"
	"github.com/study-platform/pkg/logger"
)

type LectureResourceHandler struct {
	lectureResourceService *service.LectureResourceService
	logger                 logger.Logger
}

func NewLectureResourceHandler(
	lectureResourceService *service.LectureResourceService,
	logger logger.Logger,
) *LectureResourceHandler {
	return &LectureResourceHandler{
		lectureResourceService: lectureResourceService,
		logger:                 logger,
	}
}

// CreateLectureResource creates a new resource for a lecture
// POST /api/v1/lectures/:lecture_id/resources
func (h *LectureResourceHandler) CreateLectureResource(c *gin.Context) {
	lectureIDParam := c.Param("lecture_id")
	lectureID, err := uuid.Parse(lectureIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lecture ID"})
		return
	}

	var req struct {
		FileID       string `json:"file_id" binding:"required"`
		ResourceType string `json:"resource_type" binding:"required"`
		DisplayOrder int32  `json:"display_order"`
		IsRequired   bool   `json:"is_required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	resource := &model.LectureResource{
		LectureID:    lectureID,
		FileID:       fileID,
		ResourceType: req.ResourceType,
		DisplayOrder: req.DisplayOrder,
		IsRequired:   req.IsRequired,
	}

	err = h.lectureResourceService.CreateResource(c.Request.Context(), resource)
	if err != nil {
		h.logger.Errorf("Failed to create lecture resource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Resource created successfully",
		"resource": resource,
	})
}

// GetLectureResources gets all resources for a lecture
// GET /api/v1/lectures/:lecture_id/resources
func (h *LectureResourceHandler) GetLectureResources(c *gin.Context) {
	lectureIDParam := c.Param("lecture_id")
	lectureID, err := uuid.Parse(lectureIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lecture ID"})
		return
	}

	resources, err := h.lectureResourceService.GetResourcesByLecture(c.Request.Context(), lectureID)
	if err != nil {
		h.logger.Errorf("Failed to get lecture resources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"count":     len(resources),
	})
}

// GetCourseResources gets all resources for all lectures in a course
// GET /api/v1/courses/:course_id/resources
func (h *LectureResourceHandler) GetCourseResources(c *gin.Context) {
	courseIDParam := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	resources, err := h.lectureResourceService.GetResourcesByCourse(c.Request.Context(), courseID)
	if err != nil {
		h.logger.Errorf("Failed to get course resources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
	})
}

// GetResource gets a specific resource by ID
// GET /api/v1/resources/:resource_id
func (h *LectureResourceHandler) GetResource(c *gin.Context) {
	resourceIDParam := c.Param("resource_id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	resource, err := h.lectureResourceService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		h.logger.Errorf("Failed to get resource: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resource": resource,
	})
}

// UpdateResource updates an existing lecture resource
// PUT /api/v1/resources/:resource_id
func (h *LectureResourceHandler) UpdateResource(c *gin.Context) {
	resourceIDParam := c.Param("resource_id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var req struct {
		ResourceType string `json:"resource_type"`
		DisplayOrder int32  `json:"display_order"`
		IsRequired   bool   `json:"is_required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing resource to preserve lecture_id and file_id
	existing, err := h.lectureResourceService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	resource := &model.LectureResource{
		ID:           resourceID,
		LectureID:    existing.LectureID,
		FileID:       existing.FileID,
		ResourceType: req.ResourceType,
		DisplayOrder: req.DisplayOrder,
		IsRequired:   req.IsRequired,
	}

	err = h.lectureResourceService.UpdateResource(c.Request.Context(), resource)
	if err != nil {
		h.logger.Errorf("Failed to update resource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Resource updated successfully",
		"resource": resource,
	})
}

// DeleteResource deletes a lecture resource
// DELETE /api/v1/resources/:resource_id
func (h *LectureResourceHandler) DeleteResource(c *gin.Context) {
	resourceIDParam := c.Param("resource_id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	err = h.lectureResourceService.DeleteResource(c.Request.Context(), resourceID)
	if err != nil {
		h.logger.Errorf("Failed to delete resource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resource deleted successfully",
	})
}

// BulkCreateResources creates multiple resources for a lecture
// POST /api/v1/lectures/:lecture_id/resources/bulk
func (h *LectureResourceHandler) BulkCreateResources(c *gin.Context) {
	lectureIDParam := c.Param("lecture_id")
	lectureID, err := uuid.Parse(lectureIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lecture ID"})
		return
	}

	var req struct {
		Resources []struct {
			FileID       string `json:"file_id" binding:"required"`
			ResourceType string `json:"resource_type" binding:"required"`
			DisplayOrder int32  `json:"display_order"`
			IsRequired   bool   `json:"is_required"`
		} `json:"resources" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resources := make([]model.LectureResource, len(req.Resources))
	for i, r := range req.Resources {
		fileID, err := uuid.Parse(r.FileID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID in resource " + strconv.Itoa(i)})
			return
		}

		resources[i] = model.LectureResource{
			FileID:       fileID,
			ResourceType: r.ResourceType,
			DisplayOrder: r.DisplayOrder,
			IsRequired:   r.IsRequired,
		}
	}

	err = h.lectureResourceService.BulkCreateResources(c.Request.Context(), lectureID, resources)
	if err != nil {
		h.logger.Errorf("Failed to bulk create resources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resources"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Resources created successfully",
		"count":   len(resources),
	})
}

// ReorderResources updates the display order of multiple resources
// PUT /api/v1/lectures/:lecture_id/resources/reorder
func (h *LectureResourceHandler) ReorderResources(c *gin.Context) {
	lectureIDParam := c.Param("lecture_id")
	lectureID, err := uuid.Parse(lectureIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lecture ID"})
		return
	}

	var req struct {
		Resources []struct {
			ResourceID   string `json:"resource_id" binding:"required"`
			DisplayOrder int32  `json:"display_order" binding:"required"`
		} `json:"resources" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceOrders := make([]struct {
		ResourceID   uuid.UUID `json:"resource_id"`
		DisplayOrder int32     `json:"display_order"`
	}, len(req.Resources))

	for i, r := range req.Resources {
		resourceID, err := uuid.Parse(r.ResourceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID in position " + strconv.Itoa(i)})
			return
		}

		resourceOrders[i] = struct {
			ResourceID   uuid.UUID `json:"resource_id"`
			DisplayOrder int32     `json:"display_order"`
		}{
			ResourceID:   resourceID,
			DisplayOrder: r.DisplayOrder,
		}
	}

	err = h.lectureResourceService.ReorderResources(c.Request.Context(), lectureID, resourceOrders)
	if err != nil {
		h.logger.Errorf("Failed to reorder resources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resources reordered successfully",
	})
}

// GetLectureWithResources gets a lecture with its resources populated
// GET /api/v1/lectures/:lecture_id/with-resources
func (h *LectureResourceHandler) GetLectureWithResources(c *gin.Context) {
	lectureIDParam := c.Param("lecture_id")
	lectureID, err := uuid.Parse(lectureIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lecture ID"})
		return
	}

	lecture, err := h.lectureResourceService.GetLectureWithResources(c.Request.Context(), lectureID)
	if err != nil {
		h.logger.Errorf("Failed to get lecture with resources: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Lecture not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lecture": lecture,
	})
}

// GetResourceDownloadURL generates a signed download URL for a lecture resource
// GET /api/v1/lecture-resources/{resourceId}/download-url
func (h *LectureResourceHandler) GetResourceDownloadURL(c *gin.Context) {
	resourceIDParam := c.Param("resource_id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the resource with file information
	resource, err := h.lectureResourceService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		h.logger.Errorf("Failed to get resource: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Check if user has access to the course containing this resource
	hasAccess, err := h.lectureResourceService.CheckUserCourseAccess(c.Request.Context(), userUUID, resource.LectureID)
	if err != nil {
		h.logger.Errorf("Failed to check course access: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify access"})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this resource"})
		return
	}

	// Generate signed URL using resource information
	signedURL, expiresAt, err := h.generateSignedURL(resource)
	if err != nil {
		h.logger.Errorf("Failed to generate signed URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Download URL generated successfully",
		"data": gin.H{
			"download_url": signedURL,
			"expires_at":   expiresAt,
		},
	})
}

// GetResourcePreviewURL generates a signed preview URL for a lecture resource
// GET /api/v1/lecture-resources/{resourceId}/preview-url
func (h *LectureResourceHandler) GetResourcePreviewURL(c *gin.Context) {
	resourceIDParam := c.Param("resource_id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the resource with file information
	resource, err := h.lectureResourceService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		h.logger.Errorf("Failed to get resource: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Check if user has access to the course containing this resource
	hasAccess, err := h.lectureResourceService.CheckUserCourseAccess(c.Request.Context(), userUUID, resource.LectureID)
	if err != nil {
		h.logger.Errorf("Failed to check course access: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify access"})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this resource"})
		return
	}

	// Generate signed URL using resource information (same as download for most file types)
	signedURL, expiresAt, err := h.generateSignedURL(resource)
	if err != nil {
		h.logger.Errorf("Failed to generate signed URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate preview URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Preview URL generated successfully",
		"data": gin.H{
			"preview_url": signedURL,
			"expires_at":  expiresAt,
		},
	})
}

// generateSignedURL generates a public S3 URL for unlimited access
func (h *LectureResourceHandler) generateSignedURL(resource *model.LectureResource) (string, string, error) {
	if resource.BucketName == "" || resource.ObjectKey == "" {
		return "", "", fmt.Errorf("missing bucket name or object key for resource %s", resource.ID)
	}

	// Always return direct S3 public URL for unlimited access
	// Since the bucket has public read policy, we don't need signed URLs
	directURL := fmt.Sprintf("https://%s.s3.ap-southeast-2.amazonaws.com/%s", resource.BucketName, resource.ObjectKey)

	// Set expiration to a far future date since public URLs don't expire
	expiresAt := time.Now().Add(365 * 24 * time.Hour) // 1 year from now

	h.logger.Infof("Generated public URL for resource %s: %s", resource.ID, directURL)

	return directURL, expiresAt.Format(time.RFC3339), nil
}
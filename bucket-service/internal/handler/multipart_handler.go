package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bucket-service/internal/middleware"
	"bucket-service/internal/model"
	"bucket-service/internal/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StartUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
	Bucket      string `json:"bucket"`
}

type StartUploadResponse struct {
	SessionID      uuid.UUID `json:"session_id"`
	UploadID       string    `json:"upload_id"`
	PresignedURLs  []PartURL `json:"presigned_urls"`
}

type PartURL struct {
	PartNumber int32  `json:"part_number"`
	URL        string `json:"url"`
}

type CompleteUploadRequest struct {
	Parts []CompletedPart `json:"parts" binding:"required"`
}

type CompletedPart struct {
	PartNumber int32  `json:"part_number" binding:"required"`
	ETag       string `json:"etag" binding:"required"`
}

func (h *FileHandler) StartMultipartUpload(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req StartUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate file size
	if req.Size > h.config.Upload.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File size exceeds limit of %d bytes", h.config.Upload.MaxFileSize)})
		return
	}

	// Determine file type and bucket
	fileType := model.GetFileType(req.ContentType)
	bucketName := h.s3Service.GetBucketName(fileType)

	// Generate object key
	sessionID := uuid.New()
	objectKey := h.generateObjectKey(sessionID, req.Filename, fileType)

	// Create multipart upload in S3
	metadata := map[string]string{
		"user-id":     userID.String(),
		"session-id":  sessionID.String(),
		"filename":    req.Filename,
		"upload-time": time.Now().UTC().Format(time.RFC3339),
	}

	multipartUpload, err := h.s3Service.CreateMultipartUpload(c.Request.Context(), bucketName, objectKey, req.ContentType, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create multipart upload"})
		return
	}

	// Calculate number of parts
	partSize := h.config.Upload.PartSize
	numParts := int32((req.Size + partSize - 1) / partSize) // Ceiling division

	if numParts > int32(h.config.Upload.MaxUploadParts) {
		// Abort the multipart upload
		h.s3Service.AbortMultipartUpload(c.Request.Context(), multipartUpload)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large, requires %d parts (max: %d)", numParts, h.config.Upload.MaxUploadParts)})
		return
	}

	// Generate presigned URLs for first few parts
	partNumbers := make([]int32, 0)
	maxInitialParts := int32(5) // Generate URLs for first 5 parts initially
	for i := int32(1); i <= numParts && i <= maxInitialParts; i++ {
		partNumbers = append(partNumbers, i)
	}

	expires := time.Duration(h.config.Upload.PresignedURLExpiration) * time.Second
	urls, err := h.s3Service.GetMultipartPresignedURLs(c.Request.Context(), multipartUpload, partNumbers, expires)
	if err != nil {
		// Abort the multipart upload
		h.s3Service.AbortMultipartUpload(c.Request.Context(), multipartUpload)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URLs"})
		return
	}

	// Save upload session to database
	session := &model.UploadSession{
		ID:          sessionID,
		UploadID:    multipartUpload.UploadID,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		TotalSize:   req.Size,
		BucketName:  bucketName,
		ObjectKey:   objectKey,
		UserID:      userID,
		Status:      "active",
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour), // Session expires in 24 hours
	}

	if err := h.fileRepo.CreateUploadSession(c.Request.Context(), session); err != nil {
		// Abort the multipart upload
		h.s3Service.AbortMultipartUpload(c.Request.Context(), multipartUpload)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save upload session"})
		return
	}

	// Prepare response
	presignedURLs := make([]PartURL, len(partNumbers))
	for i, partNumber := range partNumbers {
		presignedURLs[i] = PartURL{
			PartNumber: partNumber,
			URL:        urls[i],
		}
	}

	response := &StartUploadResponse{
		SessionID:     sessionID,
		UploadID:      multipartUpload.UploadID,
		PresignedURLs: presignedURLs,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *FileHandler) GetPartURLs(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get upload session
	session, err := h.fileRepo.GetUploadSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload session not found"})
		return
	}

	// Check ownership
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if session is still active
	if session.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Upload session is not active"})
		return
	}

	// Parse part numbers from query
	partNumbersStr := c.Query("part_numbers")
	if partNumbersStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Part numbers are required"})
		return
	}

	var partNumbers []int32
	for _, partStr := range strings.Split(partNumbersStr, ",") {
		partNum, err := strconv.ParseInt(strings.TrimSpace(partStr), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid part number: " + partStr})
			return
		}
		partNumbers = append(partNumbers, int32(partNum))
	}

	// Create multipart upload object for S3 service
	multipartUpload := &service.MultipartUpload{
		UploadID:   session.UploadID,
		BucketName: session.BucketName,
		ObjectKey:  session.ObjectKey,
	}

	expires := time.Duration(h.config.Upload.PresignedURLExpiration) * time.Second
	urls, err := h.s3Service.GetMultipartPresignedURLs(c.Request.Context(), multipartUpload, partNumbers, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URLs"})
		return
	}

	// Prepare response
	presignedURLs := make([]PartURL, len(partNumbers))
	for i, partNumber := range partNumbers {
		presignedURLs[i] = PartURL{
			PartNumber: partNumber,
			URL:        urls[i],
		}
	}

	c.JSON(http.StatusOK, gin.H{"presigned_urls": presignedURLs})
}

func (h *FileHandler) CompleteMultipartUpload(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get upload session
	session, err := h.fileRepo.GetUploadSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload session not found"})
		return
	}

	// Check ownership
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if session is still active
	if session.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Upload session is not active"})
		return
	}

	// Create multipart upload object for S3 service
	multipartUpload := &service.MultipartUpload{
		UploadID:   session.UploadID,
		BucketName: session.BucketName,
		ObjectKey:  session.ObjectKey,
		Parts:      make([]types.CompletedPart, len(req.Parts)),
	}

	for i, part := range req.Parts {
		multipartUpload.Parts[i] = types.CompletedPart{
			PartNumber: aws.Int32(part.PartNumber),
			ETag:       aws.String(part.ETag),
		}
	}

	// Complete multipart upload in S3
	result, err := h.s3Service.CompleteMultipartUpload(c.Request.Context(), multipartUpload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete multipart upload"})
		return
	}

	// Create file record in database
	fileID := uuid.New()
	file := &model.File{
		ID:               fileID,
		Filename:         session.Filename,
		OriginalFilename: session.Filename,
		ContentType:      session.ContentType,
		SizeBytes:        session.TotalSize,
		BucketName:       session.BucketName,
		ObjectKey:        session.ObjectKey,
		UploadUserID:     userID,
		IsPublic:         false, // Default to private for multipart uploads
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.fileRepo.CreateFile(c.Request.Context(), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	// Update session status
	session.Status = "completed"
	session.UploadedSize = session.TotalSize
	if err := h.fileRepo.UpdateUploadSession(c.Request.Context(), session); err != nil {
		// Log error but don't fail the request
	}

	response := &UploadResponse{
		FileID:      fileID,
		Filename:    session.Filename,
		Size:        session.TotalSize,
		ContentType: session.ContentType,
		URL:         *result.Location,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *FileHandler) AbortMultipartUpload(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get upload session
	session, err := h.fileRepo.GetUploadSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload session not found"})
		return
	}

	// Check ownership
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Create multipart upload object for S3 service
	multipartUpload := &service.MultipartUpload{
		UploadID:   session.UploadID,
		BucketName: session.BucketName,
		ObjectKey:  session.ObjectKey,
	}

	// Abort multipart upload in S3
	if err := h.s3Service.AbortMultipartUpload(c.Request.Context(), multipartUpload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to abort multipart upload"})
		return
	}

	// Update session status
	session.Status = "aborted"
	if err := h.fileRepo.UpdateUploadSession(c.Request.Context(), session); err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload aborted successfully"})
}

func (h *FileHandler) generateObjectKey(sessionID uuid.UUID, filename string, fileType model.FileType) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s%s", fileType, timestamp, sessionID.String(), ext)
}